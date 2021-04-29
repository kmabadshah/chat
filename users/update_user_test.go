package users

import (
	"bytes"
	"encoding/json"
	"github.com/kmabadshah/chat"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestUpdateUser(t *testing.T) {
	chat.ClearAllTables()

	router := NewRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	user := CreateTestUser(t)

	sendUpdateReq := func(t *testing.T, reqBody interface{}) *http.Response {
		encodedReqBody, err := json.Marshal(reqBody)
		chat.AssertTestErr(t, err)

		url := testServer.URL + "/users/" + strconv.Itoa(user.ID)
		req, err := http.NewRequest("PUT", url, bytes.NewReader(encodedReqBody))
		chat.AssertTestErr(t, err)

		httpClient := http.Client{}
		res, err := httpClient.Do(req)
		chat.AssertTestErr(t, err)

		return res
	}

	t.Run("update password", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"pass": "badshah",
		}
		res := sendUpdateReq(t, reqBody)

		chat.AssertTestStatusCode(t, res.StatusCode, http.StatusOK)

		resBody, err := ioutil.ReadAll(res.Body)
		chat.AssertTestErr(t, err)

		var decodedResBody map[string]interface{}
		err = json.Unmarshal(resBody, &decodedResBody)
		chat.AssertTestErr(t, err)

		if decodedResBody["pass"] != reqBody["pass"] || decodedResBody["uname"] != user.Uname {
			t.Errorf("user not updated, got %#v", decodedResBody)
		}
	})

	t.Run("change username to something that is not in use", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"uname": "adnan",
		}
		res := sendUpdateReq(t, reqBody)

		chat.AssertTestStatusCode(t, res.StatusCode, http.StatusOK)

		resBody, err := ioutil.ReadAll(res.Body)
		chat.AssertTestErr(t, err)

		var decodedResBody map[string]interface{}
		err = json.Unmarshal(resBody, &decodedResBody)
		chat.AssertTestErr(t, err)

		if decodedResBody["uname"] != reqBody["uname"] || decodedResBody["pass"] != "badshah" {
			t.Errorf("user not updated, got %#v", decodedResBody)
		}
	})

	t.Run("change username to something that is in use", func(t *testing.T) {
		user2 := CreateTestUser(t)

		reqBody := map[string]interface{}{
			"uname": user2.Uname,
		}
		res := sendUpdateReq(t, reqBody)

		chat.AssertTestStatusCode(t, res.StatusCode, http.StatusBadRequest)

		resBody, err := ioutil.ReadAll(res.Body)
		chat.AssertTestErr(t, err)

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

		chat.AssertTestStatusCode(t, res.StatusCode, http.StatusBadRequest)

		resBody, err := ioutil.ReadAll(res.Body)
		chat.AssertTestErr(t, err)

		got := string(resBody)
		want := errUpdateBody
		if got != want {
			t.Errorf("invalid resBody, wanted %#v, got %#v", want, got)
		}
	})
}
