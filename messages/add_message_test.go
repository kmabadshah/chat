package messages

import (
	"bytes"
	"encoding/json"
	"github.com/kmabadshah/chat"
	"github.com/kmabadshah/chat/users"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddMessage(t *testing.T) {
	chat.ClearAllTables()

	router := NewRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	user1 := users.CreateTestUser(t)
	user2 := users.CreateTestUser(t)

	t.Run("valid req body", func(t *testing.T) {
		url := testServer.URL + "/messages"

		reqBody := map[string]interface{}{
			"srcID": user1.ID,
			"tarID": user2.ID,
			"text":  "hello user2",
		}
		encodedReqBody, err := json.Marshal(reqBody)
		chat.AssertTestErr(t, err)

		res, err := http.Post(url, "application/json", bytes.NewReader(encodedReqBody))
		chat.AssertTestErr(t, err)

		chat.AssertTestStatusCode(t, res.StatusCode, http.StatusOK)

		resBody, err := ioutil.ReadAll(res.Body)
		chat.AssertTestErr(t, err)

		want := resAddSuccess
		got := string(resBody)

		if got != want {
			t.Errorf("invalid resBody, wanted %#v, got %#v", want, got)
		}
	})
}
