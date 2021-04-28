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

	user := createTestUser(t)

	t.Run("update password", func(t *testing.T) {
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

	t.Run("change username to something that is not in use", func(t *testing.T) {
		reqBody, err := json.Marshal(map[string]string{
			"uname": "adnan",
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

		if decodedResBody["uname"] != "adnan" || decodedResBody["pass"] != "badshah" {
			t.Errorf("user not updated, got %#v", decodedResBody)
		}
	})

	t.Run("change username to something that is in use", func(t *testing.T) {
		user2 := createTestUser(t)

		reqBody, err := json.Marshal(map[string]string{
			"uname": user2.Uname,
		})
		assertTestErr(t, err)

		url := testServer.URL + "/users/" + strconv.Itoa(user.ID)
		req, err := http.NewRequest("PUT", url, bytes.NewReader(reqBody))
		assertTestErr(t, err)

		httpClient := http.Client{}
		res, err := httpClient.Do(req)
		assertTestErr(t, err)

		assertTestStatusCode(t, res.StatusCode, http.StatusBadRequest)

		resBody, err := ioutil.ReadAll(res.Body)
		assertTestErr(t, err)

		got := string(resBody)
		want := errUnameInUse

		if got != want {
			t.Errorf("invalid resBody, wanted %#v, got %#v", want, got)
		}
	})

	t.Run("reqBody does not contain a valid uname or pass field", func(t *testing.T) {
		reqBody, err := json.Marshal(map[int]int{
			10: 20,
			20: 30,
		})
		assertTestErr(t, err)

		url := testServer.URL + "/users/" + strconv.Itoa(user.ID)
		req, err := http.NewRequest("PUT", url, bytes.NewReader(reqBody))
		assertTestErr(t, err)

		httpClient := http.Client{}
		res, err := httpClient.Do(req)
		assertTestErr(t, err)

		assertTestStatusCode(t, res.StatusCode, http.StatusBadRequest)

		resBody, err := ioutil.ReadAll(res.Body)
		assertTestErr(t, err)

		got := string(resBody)
		want := "invalid body, must contain a valid uname and/or pass field"

		if got != want {
			t.Errorf("invalid resBody, wanted %#v, got %#v", want, got)
		}
	})
}
