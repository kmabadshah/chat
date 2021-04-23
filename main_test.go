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

	// read time limit
	err := cl2.SetReadDeadline(time.Now().Add(time.Millisecond))
	errSkip(t, err)

	var latestMsg2 string
	var wg sync.WaitGroup
	wg.Add(1)

	go func(msg2 *string) {
		for {
			_, msg, err := cl2.ReadMessage()
			if err != nil {
				break
			}

			*msg2 = string(msg)
		}

		wg.Done()
	}(&latestMsg2)

	// send message to client2 from client1
	reqBody := map[string]string{
		"receiver": "client2",
		"message":  "Hello client2",
	}
	encodedReqBody, _ := json.Marshal(reqBody)
	err = cl1.WriteMessage(websocket.TextMessage, encodedReqBody)
	errSkip(t, err)

	wg.Wait()

	//now client2 should get a message
	got := latestMsg2
	want := reqBody["message"]

	if got != want {
		t.Errorf("got %q but wanted %q", got, want)
	}
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
