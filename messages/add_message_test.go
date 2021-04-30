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

	url := testServer.URL + "/messages"

	x := func(t *testing.T, reqBody map[string]interface{}, wantStatus int, wantBody string) {
		encodedReqBody, err := json.Marshal(reqBody)
		chat.AssertTestErr(t, err)

		res, err := http.Post(url, "application/json", bytes.NewReader(encodedReqBody))
		chat.AssertTestErr(t, err)

		chat.AssertTestStatusCode(t, res.StatusCode, wantStatus)

		resBody, err := ioutil.ReadAll(res.Body)
		chat.AssertTestErr(t, err)

		want := wantBody
		got := string(resBody)

		if got != want {
			t.Errorf("invalid resBody, wanted %#v, got %#v", want, got)
		}
	}

	t.Run("valid req body", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"srcID": user1.ID,
			"tarID": user2.ID,
			"text":  "hello user2",
		}
		x(t, reqBody, http.StatusOK, resAddSuccess)
	})

	t.Run("with invalid syntax", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"srcID": user1.ID,
			"tarID": user2.ID,
			"text":  "hello user2",
			"10":    20,
			"30":    40,
		}
		x(t, reqBody, http.StatusBadRequest, errReqBody)
	})

	t.Run("with valid syntax but invalid data", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"srcID": user1.ID,
			"tarID": user2.ID,
			"text":  "hello user2",
			"10":    20,
			"30":    40,
		}
		x(t, reqBody, http.StatusBadRequest, errReqBody)
	})
}
