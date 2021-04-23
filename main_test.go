package chat

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestOneClient(t *testing.T) {
	x := newX()
	testServer := httptest.NewServer(x.router)
	defer testServer.Close()

	url := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/connect/client1"

	t.Run("client added after successful connection", func(t *testing.T) {
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {

			t.Fatal("dialErr: ", err)
		}

		got := len(x.clients)
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
			got = len(x.clients)
			want = 0

			if got != want {
				t.Errorf("client not removed")
			}

		})
	})

}

func TestTwoClients(t *testing.T) {
	x := newX()
	testServer := httptest.NewServer(x.router)
	defer testServer.Close()

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

func testClientMsg(t *testing.T, tname string, clSend *websocket.Conn, reqBody map[string]string, clRecv *websocket.Conn) {
	t.Run(tname, func(t *testing.T) {
		// read time limit
		err := clRecv.SetReadDeadline(time.Now().Add(time.Millisecond))
		errSkip(t, err)

		var latestMsgRecv string
		var wgRecv sync.WaitGroup
		wgRecv.Add(1)

		go func(latestMsgRecv *string) {
			for {
				_, msg, err := clRecv.ReadMessage()
				if err != nil {
					break
				}

				*latestMsgRecv = string(msg)
			}

			wgRecv.Done()
		}(&latestMsgRecv)

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
	if err != nil {
		fmt.Println(err)
		t.Skip()
	}
}
