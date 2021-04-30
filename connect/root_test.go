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

	user3 := users.CreateTestUser(t)
	url3 := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/connect/" + strconv.Itoa(user3.ID)
	conn3, res3, err := websocket.DefaultDialer.Dial(url3, nil)
	chat.AssertTestErr(t, err)
	chat.AssertTestStatusCode(t, res3.StatusCode, http.StatusSwitchingProtocols)

	t.Run("user2 gets notified when user1 sends new message", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"type": "message",
			"uid":  user2.ID,
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
		encodedReqBody, err := json.Marshal(reqBody)
		chat.AssertTestErr(t, err)

		err = conn3.WriteMessage(websocket.TextMessage, encodedReqBody)
		chat.AssertTestErr(t, err)

		var decodedData1 map[string]interface{}
		_, encodedData1, err := conn1.ReadMessage()
		err = json.Unmarshal(encodedData1, &decodedData1)
		chat.AssertTestErr(t, err)

		var decodedData2 map[string]interface{}
		_, encodedData2, err := conn2.ReadMessage()
		err = json.Unmarshal(encodedData2, &decodedData2)
		chat.AssertTestErr(t, err)

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
