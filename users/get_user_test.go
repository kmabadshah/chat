package users

import (
	"encoding/json"
	"github.com/kmabadshah/chat"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGetUser(t *testing.T) {
	chat.ClearAllTables()

	user := createTestUser(t)

	router := NewRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	t.Run("valid id", func(t *testing.T) {
		url := testServer.URL + "/users/" + strconv.Itoa(user.ID)
		res, err := http.Get(url)
		chat.AssertTestErr(t, err)

		t.Run("check code", func(t *testing.T) {
			chat.AssertTestStatusCode(t, res.StatusCode, http.StatusOK)
		})

		t.Run("check res body", func(t *testing.T) {
			resBody, err := ioutil.ReadAll(res.Body)
			chat.AssertTestErr(t, err)

			var decodedResBody map[string]interface{}
			err = json.Unmarshal(resBody, &decodedResBody)
			chat.AssertTestErr(t, err)

			if decodedResBody["uname"] != user.Uname ||
				decodedResBody["pass"] != user.Pass {

				t.Errorf("did not get two users, got %#v", decodedResBody)
			}
		})
	})

	t.Run("invalid id", func(t *testing.T) {
		url := testServer.URL + "/users/" + strconv.Itoa(-1)
		res, err := http.Get(url)
		chat.AssertTestErr(t, err)

		t.Run("check code", func(t *testing.T) {
			chat.AssertTestStatusCode(t, res.StatusCode, http.StatusNotFound)
		})
	})
}
