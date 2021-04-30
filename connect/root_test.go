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

	user1 := users.CreateTestUser(t)
	url1 := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/connect/" + strconv.Itoa(user1.ID)
	conn1, res1, err := websocket.DefaultDialer.Dial(url1, nil)
	chat.AssertTestErr(t, err)
	chat.AssertTestStatusCode(t, res1.StatusCode, http.StatusSwitchingProtocols)

	user2 := users.CreateTestUser(t)
	url2 := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/connect/" + strconv.Itoa(user2.ID)
	conn2, res2, err := websocket.DefaultDialer.Dial(url2, nil)
	chat.AssertTestErr(t, err)
	chat.AssertTestStatusCode(t, res2.StatusCode, http.StatusSwitchingProtocols)

	t.Run("user2 gets notified when user1 sends new message", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"type":  "message",
			"tarID": user2.ID,
		}
		encodedReqBody, err := json.Marshal(reqBody)
		chat.AssertTestErr(t, err)

		err = conn1.WriteMessage(websocket.TextMessage, encodedReqBody)
		chat.AssertTestErr(t, err)

		var decodedData map[string]interface{}
		_, encodedData, err := conn2.ReadMessage()
		err = json.Unmarshal(encodedData, &decodedData)
		chat.AssertTestErr(t, err)

		want := map[string]interface{}{
			"type":  "message",
			"srcID": float64(user1.ID),
		}

		if !reflect.DeepEqual(want, decodedData) {
			t.Errorf("got %#v, wanted %#v", decodedData, want)
		}
	})
}
