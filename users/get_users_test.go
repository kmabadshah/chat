package users

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUsers(t *testing.T) {
	clearUserTable()

	user1 := User{
		Uname: "adnan1",
		Pass:  "badshah1",
	}
	user2 := User{
		Uname: "adnan2",
		Pass:  "badshah2",
	}
	_, err := db.Model(&user1).Insert()
	assertTestErr(t, err)
	_, err = db.Model(&user2).Insert()
	assertTestErr(t, err)

	router := newRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	url := testServer.URL + "/users"
	res, err := http.Get(url)
	assertTestErr(t, err)

	resBody, err := ioutil.ReadAll(res.Body)
	assertTestErr(t, err)

	t.Run("check code", func(t *testing.T) {
		assertTestStatusCode(t, res.StatusCode, http.StatusOK)
	})

	t.Run("check body", func(t *testing.T) {
		var decodedResBody []map[string]interface{}
		err = json.Unmarshal(resBody, &decodedResBody)
		assertTestErr(t, err)

		if len(decodedResBody) != 2 ||
			decodedResBody[0]["uname"] != user1.Uname ||
			decodedResBody[1]["uname"] != user2.Uname {

			t.Errorf("did not get the two created users")
		}
	})
}
