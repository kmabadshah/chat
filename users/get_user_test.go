package users

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGetUser(t *testing.T) {
	clearUserTable()

	user := User{
		Uname: "adnan1",
		Pass:  "badshah1",
	}
	_, err := db.Model(&user).Insert()
	assertTestErr(t, err)

	router := newRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	t.Run("valid id", func(t *testing.T) {
		url := testServer.URL + "/users/" + strconv.Itoa(user.ID)
		res, err := http.Get(url)
		assertTestErr(t, err)

		t.Run("check code", func(t *testing.T) {
			assertTestStatusCode(t, res.StatusCode, http.StatusOK)
		})

		t.Run("check res body", func(t *testing.T) {
			resBody, err := ioutil.ReadAll(res.Body)
			assertTestErr(t, err)

			var decodedResBody map[string]interface{}
			err = json.Unmarshal(resBody, &decodedResBody)
			assertTestErr(t, err)

			if decodedResBody["uname"] != user.Uname ||
				decodedResBody["pass"] != user.Pass {

				t.Errorf("did not get two users, got %#v", decodedResBody)
			}
		})
	})

	t.Run("invalid id", func(t *testing.T) {
		url := testServer.URL + "/users/" + strconv.Itoa(-1)
		res, err := http.Get(url)
		assertTestErr(t, err)

		t.Run("check code", func(t *testing.T) {
			assertTestStatusCode(t, res.StatusCode, http.StatusNotFound)
		})
	})
}
