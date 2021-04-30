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

	sendAddFriendReq := func(t *testing.T, reqBody interface{}, wantStatus int, wantBody string) {
		encodedReqBody, err := json.Marshal(reqBody)
		chat.AssertTestErr(t, err)

		res, err := http.Post(url, "application/json", bytes.NewReader(encodedReqBody))
		chat.AssertTestErr(t, err)

		chat.AssertTestStatusCode(t, res.StatusCode, wantStatus)

		if wantBody != "" {
			resBody, err := ioutil.ReadAll(res.Body)
			chat.AssertTestErr(t, err)

			gotBody := string(resBody)
			if gotBody != wantBody {
				t.Errorf("invalid resp body, wanted %#v, gotBody %#v", wantBody, gotBody)
			}
		}
	}

	t.Run("valid req body", func(t *testing.T) {
		reqBody := map[string]int{
			"srcID": user1.ID,
			"tarID": user2.ID,
		}

		sendAddFriendReq(t, reqBody, http.StatusOK, "")
	})

	t.Run("invalid body syntax", func(t *testing.T) {
		reqBody := map[int]int{
			10: 20,
			30: 40,
		}

		sendAddFriendReq(t, reqBody, http.StatusBadRequest, errReqBody)
	})

	t.Run("invalid srcID field", func(t *testing.T) {
		reqBody := map[string]int{
			"srcID": -1,
			"tarID": -2,
		}

		sendAddFriendReq(t, reqBody, http.StatusBadRequest, errReqBody)
	})
}
