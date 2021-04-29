package friends

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

func TestAddFriend(t *testing.T) {
	chat.ClearAllTables()

	router := NewRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	user1 := users.CreateTestUser(t)
	user2 := users.CreateTestUser(t)

	url := testServer.URL + "/friends"

	t.Run("valid req body", func(t *testing.T) {
		reqBody := map[string]int{
			"srcID": user1.ID,
			"tarID": user2.ID,
		}

		encodedReqBody, err := json.Marshal(reqBody)
		chat.AssertTestErr(t, err)

		res, err := http.Post(url, "application/json", bytes.NewReader(encodedReqBody))
		chat.AssertTestErr(t, err)

		chat.AssertTestStatusCode(t, res.StatusCode, http.StatusOK)
	})

	t.Run("invalid body syntax", func(t *testing.T) {
		reqBody := map[int]int{
			10: 20,
			30: 40,
		}

		encodedReqBody, err := json.Marshal(reqBody)
		chat.AssertTestErr(t, err)

		res, err := http.Post(url, "application/json", bytes.NewReader(encodedReqBody))
		chat.AssertTestErr(t, err)

		chat.AssertTestStatusCode(t, res.StatusCode, http.StatusBadRequest)

		resBody, err := ioutil.ReadAll(res.Body)
		chat.AssertTestErr(t, err)

		got := string(resBody)
		want := errReqBody

		if got != want {
			t.Errorf("invalid resp body, wanted %#v, got %#v", want, got)
		}
	})
}
