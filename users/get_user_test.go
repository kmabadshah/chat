package users

import (
	"encoding/json"
	"github.com/kmabadshah/chat/shared"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

func TestGetUser(t *testing.T) {
	shared.ClearAllTables()

	user := CreateTestUser(t)

	router := NewRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	t.Run("valid uname", func(t *testing.T) {
		url := testServer.URL + "/users/" + user.Uname
		res, err := http.Get(url)
		shared.AssertTestErr(t, err)

		t.Run("check code", func(t *testing.T) {
			shared.AssertTestStatusCode(t, res.StatusCode, http.StatusOK)
		})

		t.Run("check res body", func(t *testing.T) {
			resBody, err := ioutil.ReadAll(res.Body)
			shared.AssertTestErr(t, err)

			var decodedResBody map[string]interface{}
			err = json.Unmarshal(resBody, &decodedResBody)
			shared.AssertTestErr(t, err)

			if reflect.TypeOf(decodedResBody["id"]).String() == "float64" {
				decodedResBody["id"] = int(decodedResBody["id"].(float64))
			}

			if decodedResBody["id"] != nil {
				decodedResBody["id"] = decodedResBody["id"].(int)
			}

			var want map[string]interface{}
			err = mapstructure.Decode(user, &want)
			shared.AssertTestErr(t, err)

			if !reflect.DeepEqual(want, decodedResBody) {
				t.Errorf("wanted %#v, got %#v", want, decodedResBody)
			}
		})
	})

	t.Run("invalid name", func(t *testing.T) {
		url := testServer.URL + "/users/" + strconv.Itoa(-1)
		res, err := http.Get(url)
		shared.AssertTestErr(t, err)

		t.Run("check status code", func(t *testing.T) {
			shared.AssertTestStatusCode(t, res.StatusCode, http.StatusNotFound)
		})
	})
}
