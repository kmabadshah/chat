package friends

import (
	"bytes"
	"encoding/json"
	"github.com/kmabadshah/chat"
	"github.com/kmabadshah/chat/users"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddFriend(t *testing.T) {
	chat.ClearAllTables()

	router := NewRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	t.Run("valid id and req body", func(t *testing.T) {
		user1 := users.CreateTestUser(t)
		user2 := users.CreateTestUser(t)
		url := testServer.URL + "/friends"
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

	t.Run("invalid id, valid body", func(t *testing.T) {

	})
}
