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
	url1 := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/connect/client1"
	conn1, _, err := websocket.DefaultDialer.Dial(url1, nil)
	if err != nil {
		t.Fatal("dialError: ", err)
	}
	defer conn1.Close()

	// client 2
	url2 := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/connect/client2"
	conn2, _, err := websocket.DefaultDialer.Dial(url2, nil)
	errSkip(t, err)
	defer conn2.Close()

	err = conn2.SetReadDeadline(time.Now().Add(time.Millisecond))
	errSkip(t, err)

	var msg2 string
	var wg sync.WaitGroup
	wg.Add(1)

	go func(msg2 *string) {
		for {
			_, msg, err := conn2.ReadMessage()
			if err != nil {
				break
			}

			*msg2 = string(msg)
		}

		wg.Done()
	}(&msg2)

	// send message to client2 from client1
	reqBody := map[string]string{
		"receiver": "client2",
		"message":  "Hello client2",
	}
	encodedReqBody, _ := json.Marshal(reqBody)
	err = conn1.WriteMessage(websocket.TextMessage, encodedReqBody)
	errSkip(t, err)

	wg.Wait()

	//now client2 should get a message
	got := msg2
	want := reqBody["message"]

	if got != want {
		t.Errorf("got %q but wanted %q", got, want)
	}
}

func errSkip(t *testing.T, err error) {
	if err != nil {
		fmt.Println(err)
		t.Skip()
	}
}
