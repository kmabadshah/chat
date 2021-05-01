package connect

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/kmabadshah/chat"
	"github.com/kmabadshah/chat/users"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestConnectUser(t *testing.T) {
	chat.ClearAllTables()

	router := NewRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	user1, conn1 := testCreateAndConnect(t, testServer.URL)

	var user2 users.User
	var conn2 *websocket.Conn

	var user3 users.User
	var conn3 *websocket.Conn

	t.Run("when user2 joins, user1 gets notified", func(t *testing.T) {
		user2, conn2 = testCreateAndConnect(t, testServer.URL)

		want := map[string]interface{}{
			"type": "connect",
			"uid":  float64(user2.ID),
		}

		testReadMsgAndComp(t, conn1, want)
	})

	t.Run("when user3 joins", func(t *testing.T) {
		user3, conn3 = testCreateAndConnect(t, testServer.URL)

		t.Run("user1 gets notified", func(t *testing.T) {
			want := map[string]interface{}{
				"type": "connect",
				"uid":  float64(user3.ID),
			}

			testReadMsgAndComp(t, conn1, want)
		})

		t.Run("user2 gets notified", func(t *testing.T) {
			want := map[string]interface{}{
				"type": "connect",
				"uid":  float64(user3.ID),
			}

			testReadMsgAndComp(t, conn2, want)
		})
	})

	t.Run("user2 gets notified when user1 sends new message", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"type": "message",
			"uid":  user2.ID,
		}
		testSendMsg(t, reqBody, conn1)

		decodedData := testRecvMsg(t, conn2)

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
		testSendMsg(t, reqBody, conn3)

		decodedData1 := testRecvMsg(t, conn1)
		decodedData2 := testRecvMsg(t, conn2)

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

	t.Run("user2 gets friend request from user1", func(t *testing.T) {
		encMsg, err := json.Marshal(map[string]interface{}{
			"type": "friend",
			"uid":  float64(user2.ID),
		})
		chat.AssertTestErr(t, err)

		err = conn1.WriteMessage(websocket.TextMessage, encMsg)
		chat.AssertTestErr(t, err)

		want := map[string]interface{}{
			"type": "friend",
			"uid":  float64(user1.ID),
		}
		testReadMsgAndComp(t, conn2, want)
	})

	t.Run("when user3 leaves", func(t *testing.T) {
		err := conn3.Close()
		chat.AssertTestErr(t, err)

		t.Run("user2 gets notified", func(t *testing.T) {
			want := map[string]interface{}{
				"type": "leave",
				"uid":  float64(user3.ID),
			}

			testReadMsgAndComp(t, conn2, want)
		})

		t.Run("user1 gets notified", func(t *testing.T) {
			want := map[string]interface{}{
				"type": "leave",
				"uid":  float64(user3.ID),
			}

			testReadMsgAndComp(t, conn1, want)
		})
	})
}
