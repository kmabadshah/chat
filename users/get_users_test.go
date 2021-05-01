package users

import (
	"encoding/json"
	"github.com/kmabadshah/chat/shared"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUsers(t *testing.T) {
	shared.ClearAllTables()

	user1 := CreateTestUser(t)
	user2 := CreateTestUser(t)

	router := NewRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	url := testServer.URL + "/users"
	res, err := http.Get(url)
	shared.AssertTestErr(t, err)

	resBody, err := ioutil.ReadAll(res.Body)
	shared.AssertTestErr(t, err)

	t.Run("check code", func(t *testing.T) {
		shared.AssertTestStatusCode(t, res.StatusCode, http.StatusOK)
	})

	t.Run("check body", func(t *testing.T) {
		var decodedResBody []map[string]interface{}
		err = json.Unmarshal(resBody, &decodedResBody)
		shared.AssertTestErr(t, err)

		if len(decodedResBody) != 2 ||
			decodedResBody[0]["uname"] != user1.Uname ||
			decodedResBody[1]["uname"] != user2.Uname {

			t.Errorf("did not get the two created users")
		}
	})
}
