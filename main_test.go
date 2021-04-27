package chat

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddUser(t *testing.T) {
	_, err := db.Exec(`delete from users`)
	assertTestErr(t, err)

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

	t.Run("status code and body", func(t *testing.T) {
		gotStatus := res.StatusCode
		wantStatus := http.StatusOK

		if gotStatus != wantStatus {
			t.Errorf("invalid code, wanted %#v, got %#v", wantStatus, gotStatus)
		}

		for k, v := range reqBody {
			if decodedResBody[k] != v {
				t.Fatalf("invalid resp, got %#v, wanted %#v", decodedResBody, reqBody)
			}
		}
	})
}

func assertTestErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
