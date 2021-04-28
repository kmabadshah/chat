package users

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestUpdateUser(t *testing.T) {
	clearUserTable()

	router := newRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	t.Run("update password", func(t *testing.T) {
		user := createTestUser(t)

		reqBody, err := json.Marshal(map[string]string{
			"pass": "badshah",
		})
		assertTestErr(t, err)

		url := testServer.URL + "/users/" + strconv.Itoa(user.ID)
		req, err := http.NewRequest("PUT", url, bytes.NewReader(reqBody))
		assertTestErr(t, err)

		httpClient := http.Client{}
		res, err := httpClient.Do(req)
		assertTestErr(t, err)

		assertTestStatusCode(t, res.StatusCode, http.StatusOK)

		resBody, err := ioutil.ReadAll(res.Body)
		assertTestErr(t, err)
		var decodedResBody map[string]interface{}
		err = json.Unmarshal(resBody, &decodedResBody)
		assertTestErr(t, err)

		if decodedResBody["pass"] != "badshah" || decodedResBody["uname"] != user.Uname {
			t.Errorf("user not updated, got %#v", decodedResBody)
		}
	})
}
