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

	sendUpdateReq := func(t *testing.T, reqBody interface{}) *http.Response {
		encodedReqBody, err := json.Marshal(reqBody)
		assertTestErr(t, err)

		url := testServer.URL + "/users/" + strconv.Itoa(user.ID)
		req, err := http.NewRequest("PUT", url, bytes.NewReader(encodedReqBody))
		assertTestErr(t, err)

		httpClient := http.Client{}
		res, err := httpClient.Do(req)
		assertTestErr(t, err)

		return res
	}

	t.Run("update password", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"pass": "badshah",
		}
		res := sendUpdateReq(t, reqBody)

		assertTestStatusCode(t, res.StatusCode, http.StatusOK)

		resBody, err := ioutil.ReadAll(res.Body)
		assertTestErr(t, err)

		var decodedResBody map[string]interface{}
		err = json.Unmarshal(resBody, &decodedResBody)
		assertTestErr(t, err)

		if decodedResBody["pass"] != reqBody["pass"] || decodedResBody["uname"] != user.Uname {
			t.Errorf("user not updated, got %#v", decodedResBody)
		}
	})

	t.Run("change username to something that is not in use", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"uname": "adnan",
		}
		res := sendUpdateReq(t, reqBody)

		assertTestStatusCode(t, res.StatusCode, http.StatusOK)

		resBody, err := ioutil.ReadAll(res.Body)
		assertTestErr(t, err)

		var decodedResBody map[string]interface{}
		err = json.Unmarshal(resBody, &decodedResBody)
		assertTestErr(t, err)

		if decodedResBody["uname"] != reqBody["uname"] || decodedResBody["pass"] != "badshah" {
			t.Errorf("user not updated, got %#v", decodedResBody)
		}
	})

	t.Run("change username to something that is in use", func(t *testing.T) {
		user2 := createTestUser(t)

		reqBody := map[string]interface{}{
			"uname": user2.Uname,
		}
		res := sendUpdateReq(t, reqBody)

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
		reqBody := map[int]int{
			10: 20,
			20: 30,
		}
		res := sendUpdateReq(t, reqBody)

		assertTestStatusCode(t, res.StatusCode, http.StatusBadRequest)

		resBody, err := ioutil.ReadAll(res.Body)
		assertTestErr(t, err)

		got := string(resBody)
		want := errUpdateBody
		if got != want {
			t.Errorf("invalid resBody, wanted %#v, got %#v", want, got)
		}
	})
}
