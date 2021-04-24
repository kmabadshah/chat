package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var hub *Hub
var testServer *httptest.Server

func TestMain(m *testing.M) {
	hub = newHub()
	testServer = httptest.NewServer(hub.router)

	code := m.Run()
	testServer.Close()

	os.Exit(code)
}

func TestOneClient(t *testing.T) {
	url := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/connect/client1"

	t.Run("client added after successful connection", func(t *testing.T) {
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatal("dialErr: ", err)
		}

		got := len(hub.clients)
		want := 1

		if got != want {
			t.Fatal("client not registered")
		}

		t.Run("client got notified about successful connection", func(t *testing.T) {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				t.Fatal("readMsg: ", err)
			}
			got := string(msg)
			want := "connection successful"

			if got != want {
				t.Errorf("didn't get proper message, wanted %q, got %q", want, got)
			}

		})

		t.Run("client removed after connection close", func(t *testing.T) {
			err = conn.Close()
			if err != nil {
				t.Fatal("connCloseErr: ", err)
			}

			time.Sleep(time.Millisecond)
			got = len(hub.clients)
			want = 0

			if got != want {
				t.Errorf("client not removed")
			}

		})
	})

}

func TestTwoClients(t *testing.T) {
	// client 1
	cl1 := createTestClient(t, testServer, "client1")
	defer cl1.Close()

	// client 2
	cl2 := createTestClient(t, testServer, "client2")
	defer cl2.Close()

	testClientMsg(t, "client1msg -> client2", cl1, map[string]string{
		"receiver": "client2",
		"message":  "Hello client2",
	}, cl2)

	testClientMsg(t, "client2msg -> client1", cl2, map[string]string{
		"receiver": "client1",
		"message":  "Hello client1",
	}, cl1)
}

func TestUsersEndpoint(t *testing.T) {
	clearDB(t)
	defer clearDB(t)

	hub := newHub()
	testServer := httptest.NewServer(hub.router)
	defer testServer.Close()

	reqBody := map[string]string{
		"uname": "adnan",
		"pass":  "badshah",
	}

	t.Run("create a user", func(t *testing.T) {
		url := testServer.URL + "/users/"

		// POST the user
		encodedReqBody, _ := json.Marshal(reqBody)
		res, err := http.Post(url, "application/json", bytes.NewReader(encodedReqBody))
		errSkip(t, err)

		// get the response
		resBody, err := ioutil.ReadAll(res.Body)
		errSkip(t, err)

		var decodedResBody map[string]interface{}
		err = json.Unmarshal(resBody, &decodedResBody)
		errSkip(t, err)

		// check the response
		if decodedResBody["uname"] != reqBody["uname"] || decodedResBody["pass"] != reqBody["pass"] {
			t.Errorf("invalid res body, got %#v, wanted %#v", decodedResBody, reqBody)
		}
	})

	t.Run("get a user", func(t *testing.T) {
		url := testServer.URL + "/users/" + reqBody["uname"]

		res, err := http.Get(url)
		if err != nil {
			t.Fatal(err)
		}

		// decode
		resBody, err := ioutil.ReadAll(res.Body)
		errSkip(t, err)

		var decodedResBody map[string]interface{}
		err = json.Unmarshal(resBody, &decodedResBody)
		errSkip(t, err)

		if decodedResBody["uname"] != reqBody["uname"] || decodedResBody["pass"] != reqBody["pass"] {
			t.Errorf("did not get the user created earlier")
		}
	})
}

func testClientMsg(t *testing.T, tname string, clSend *websocket.Conn, reqBody map[string]string, clRecv *websocket.Conn) {
	clearDB(t)
	defer clearDB(t)

	t.Run(tname, func(t *testing.T) {
		// read time limit
		err := clRecv.SetReadDeadline(time.Now().Add(time.Millisecond))
		errSkip(t, err)

		var latestMsgRecv string
		var wgRecv sync.WaitGroup
		wgRecv.Add(1)

		go func() {
			for {
				_, msg, err := clRecv.ReadMessage()
				if err != nil {
					break
				}

				latestMsgRecv = string(msg)
			}

			wgRecv.Done()
		}()

		for i := 0; i < 10; i++ {
			// send message to clRecv from clSend
			encodedReqBody, _ := json.Marshal(reqBody)
			err = clSend.WriteMessage(websocket.TextMessage, encodedReqBody)
			errSkip(t, err)

			wgRecv.Wait()

			//now clientRecv should get a message
			got := latestMsgRecv
			want := reqBody["message"]

			if got != want {
				t.Errorf("got %q but wanted %q", got, want)
			}
		}
	})
}

func createTestClient(t *testing.T, testServer *httptest.Server, cname string) *websocket.Conn {
	url := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/connect/" + cname
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal("dialError: ", err)
	}

	return conn
}

func errSkip(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		fmt.Println(err)
		t.Skip()
	}
}

func clearDB(t *testing.T) {
	t.Helper()
	// clean the user table
	_, err := db.Model((*User)(nil)).Exec(`
		delete from users
	`)
	if err != nil {
		t.Fatal(err)
	}
}
