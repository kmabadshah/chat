package users

import (
	"bytes"
	"encoding/json"
	"github.com/kmabadshah/chat/shared"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddUser(t *testing.T) {
	shared.ClearAllTables()

	router := NewRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	url := testServer.URL + "/users"
	reqBody := map[string]string{
		"uname": "adnan",
		"pass":  "badshah",
	}
	encodedReqBody, err := json.Marshal(reqBody)
	shared.AssertTestErr(t, err)

	res, err := http.Post(url, "application/json", bytes.NewReader(encodedReqBody))
	shared.AssertTestErr(t, err)

	resBody, err := ioutil.ReadAll(res.Body)
	shared.AssertTestErr(t, err)

	var decodedResBody map[string]interface{}
	err = json.Unmarshal(resBody, &decodedResBody)
	shared.AssertTestErr(t, err)

	t.Run("check status code", func(t *testing.T) {
		shared.AssertTestStatusCode(t, res.StatusCode, http.StatusOK)
	})
	t.Run("check res body", func(t *testing.T) {
		for k, v := range reqBody {
			if decodedResBody[k] != v {
				t.Errorf("invalid resp, got %#v, wanted %#v", decodedResBody, reqBody)
			}
		}
	})

	t.Run("check if user stored in db", func(t *testing.T) {
		var users []User
		err := shared.DB.Model(&users).Select()
		shared.AssertTestErr(t, err)

		if len(users) != 1 || users[0].Uname != reqBody["uname"] {
			t.Errorf("user not stored in db")
		}
	})

	t.Run("invalid req body", func(t *testing.T) {
		reqBody := map[string]string{
			"pass": "badshah",
		}
		encodedReqBody, err := json.Marshal(reqBody)
		shared.AssertTestErr(t, err)

		res, err := http.Post(url, "application/json", bytes.NewReader(encodedReqBody))
		shared.AssertTestErr(t, err)

		gotStatus := res.StatusCode
		wantStatus := http.StatusBadRequest
		if gotStatus != wantStatus {
			t.Errorf("invalid code, wanted %#v, got %#v", wantStatus, gotStatus)
		}
	})

	t.Run("storing user that already exists", func(t *testing.T) {
		url := testServer.URL + "/users"
		reqBody := map[string]string{
			"uname": "adnan",
			"pass":  "badshah",
		}
		encodedReqBody, err := json.Marshal(reqBody)
		shared.AssertTestErr(t, err)

		res, err := http.Post(url, "application/json", bytes.NewReader(encodedReqBody))
		shared.AssertTestErr(t, err)

		shared.AssertTestStatusCode(t, res.StatusCode, http.StatusBadRequest)
	})
}
