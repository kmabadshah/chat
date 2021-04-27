package users

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddUser(t *testing.T) {
	clearUserTable()

	router := newRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	url := testServer.URL + "/users"
	reqBody := map[string]string{
		"uname": "adnan",
		"pass":  "badshah",
	}
	encodedReqBody, err := json.Marshal(reqBody)
	assertTestErr(t, err)

	res, err := http.Post(url, "application/json", bytes.NewReader(encodedReqBody))
	assertTestErr(t, err)

	resBody, err := ioutil.ReadAll(res.Body)
	assertTestErr(t, err)

	var decodedResBody map[string]interface{}
	err = json.Unmarshal(resBody, &decodedResBody)
	assertTestErr(t, err)

	t.Run("check status code", func(t *testing.T) {
		assertTestStatusCode(t, res.StatusCode, http.StatusOK)
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
		err := db.Model(&users).Select()
		assertTestErr(t, err)

		if len(users) != 1 || users[0].Uname != reqBody["uname"] {
			t.Errorf("user not stored in db")
		}
	})

	t.Run("invalid req body", func(t *testing.T) {
		reqBody := map[string]string{
			"pass": "badshah",
		}
		encodedReqBody, err := json.Marshal(reqBody)
		assertTestErr(t, err)

		res, err := http.Post(url, "application/json", bytes.NewReader(encodedReqBody))
		assertTestErr(t, err)

		gotStatus := res.StatusCode
		wantStatus := http.StatusBadRequest
		if gotStatus != wantStatus {
			t.Errorf("invalid code, wanted %#v, got %#v", wantStatus, gotStatus)
		}
	})
}
