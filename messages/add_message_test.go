package messages

import (
	"bytes"
	"encoding/json"
	"github.com/kmabadshah/chat"
	"github.com/kmabadshah/chat/users"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddMessage(t *testing.T) {
	chat.ClearAllTables()

	router := NewRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	user1 := users.CreateTestUser(t)
	user2 := users.CreateTestUser(t)

	url := testServer.URL + "/messages"

	tests := []struct {
		name       string
		reqBody    interface{}
		wantStatus int
		wantRes    string
	}{
		{
			"valid req body",
			map[string]interface{}{
				"srcID": user1.ID,
				"tarID": user2.ID,
				"text":  "hello user2",
			},
			http.StatusOK,
			resAddSuccess,
		},

		{
			"invalid reqBody syntax",
			map[string]interface{}{
				"srcID": user1.ID,
				"tarID": user2.ID,
				"text":  "hello user2",
				"10":    20,
				"30":    40,
			},
			http.StatusBadRequest,
			errReqBody,
		},

		{
			"valid syntax but invalid data",
			map[string]interface{}{
				"srcID": -1,
				"tarID": -2,
				"text":  "hello world",
			},
			http.StatusBadRequest,
			errReqBody,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Helper()
			encodedReqBody, err := json.Marshal(test.reqBody)
			chat.AssertTestErr(t, err)

			res, err := http.Post(url, "application/json", bytes.NewReader(encodedReqBody))
			chat.AssertTestErr(t, err)

			chat.AssertTestStatusCode(t, res.StatusCode, test.wantStatus)

			resBody, err := ioutil.ReadAll(res.Body)
			chat.AssertTestErr(t, err)

			want := test.wantRes
			got := string(resBody)

			if got != want {
				t.Errorf("invalid resBody, wanted %#v, got %#v", want, got)
			}
		})
	}
}
