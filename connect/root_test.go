package connect

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/kmabadshah/chat"
	"github.com/kmabadshah/chat/users"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestConnectUser(t *testing.T) {
	chat.ClearAllTables()

	router := NewRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	user1, conn1 := createAndConnect(t, testServer.URL)

	var user2 users.User
	var conn2 *websocket.Conn

	var user3 users.User
	var conn3 *websocket.Conn

	t.Run("when user2 joins, user1 gets notified", func(t *testing.T) {
		user2, conn2 = createAndConnect(t, testServer.URL)

		want := map[string]interface{}{
			"type": "connect",
			"id":   float64(user2.ID),
		}

		connAndComp(t, conn1, want)
	})

	t.Run("when user3 joins", func(t *testing.T) {
		user3, conn3 = createAndConnect(t, testServer.URL)

		t.Run("user1 gets notified", func(t *testing.T) {
			want := map[string]interface{}{
				"type": "connect",
				"id":   float64(user3.ID),
			}

			connAndComp(t, conn1, want)
		})

		t.Run("user2 gets notified", func(t *testing.T) {
			want := map[string]interface{}{
				"type": "connect",
				"id":   float64(user3.ID),
			}

			connAndComp(t, conn2, want)
		})
	})

	t.Run("user2 gets notified when user1 sends new message", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"type": "message",
			"uid":  user2.ID,
		}
		sendMsg(t, reqBody, conn1)

		decodedData := recvMsg(t, conn2)

		want := map[string]interface{}{
			"type": "message",
			"uid":  float64(user1.ID),
		}

		if !reflect.DeepEqual(want, decodedData) {
			t.Errorf("got %#v, wanted %#v", decodedData, want)
		}
	})

	t.Run("user2 and user1 notified when user3 changes uname", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"type": "broadcast",
		}
		sendMsg(t, reqBody, conn3)

		decodedData1 := recvMsg(t, conn1)
		decodedData2 := recvMsg(t, conn2)

		want := map[string]interface{}{
			"type": "user",
			"uid":  float64(user3.ID),
		}

		if !reflect.DeepEqual(want, decodedData1) {
			t.Errorf("got %#v, wanted %#v", decodedData1, want)
		}
		if !reflect.DeepEqual(want, decodedData2) {
			t.Errorf("got %#v, wanted %#v", decodedData2, want)
		}
	})
}

func createAndConnect(t *testing.T, url string) (users.User, *websocket.Conn) {
	t.Helper()

	user := users.CreateTestUser(t)
	url1 := "ws" + strings.TrimPrefix(url, "http") + "/connect/" + strconv.Itoa(user.ID)
	conn, res, err := websocket.DefaultDialer.Dial(url1, nil)
	chat.AssertTestErr(t, err)
	chat.AssertTestStatusCode(t, res.StatusCode, http.StatusSwitchingProtocols)

	return user, conn
}

func sendMsg(t *testing.T, reqBody map[string]interface{}, conn *websocket.Conn) {
	t.Helper()

	encodedReqBody, err := json.Marshal(reqBody)
	chat.AssertTestErr(t, err)

	err = conn.WriteMessage(websocket.TextMessage, encodedReqBody)
	chat.AssertTestErr(t, err)
}

func recvMsg(t *testing.T, conn *websocket.Conn) map[string]interface{} {
	t.Helper()

	_, encodedData, err := conn.ReadMessage()
	var decodedData map[string]interface{}
	err = json.Unmarshal(encodedData, &decodedData)
	chat.AssertTestErr(t, err)

	return decodedData
}

func connAndComp(t *testing.T, conn *websocket.Conn, want map[string]interface{}) {
	t.Helper()

	_, encData, err := conn.ReadMessage()
	chat.AssertTestErr(t, err)

	var decData map[string]interface{}
	err = json.Unmarshal(encData, &decData)
	chat.AssertTestErr(t, err)

	if !reflect.DeepEqual(decData, want) {
		t.Errorf("got %#v, wanted %#v", decData, want)
	}
}
