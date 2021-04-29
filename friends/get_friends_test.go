package friends

import (
	"encoding/json"
	"github.com/kmabadshah/chat"
	"github.com/kmabadshah/chat/users"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

func TestGetFriends(t *testing.T) {
	chat.ClearAllTables()

	router := NewRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	user1 := users.CreateTestUser(t)
	user2 := users.CreateTestUser(t)

	// add a friend
	friend := Friend{SrcID: user1.ID, TarID: user2.ID}
	result, err := chat.DB.Model(&friend).Insert()
	if err != nil || result.RowsAffected() != 1 {
		t.Fatal(err)
	}

	t.Run("valid id", func(t *testing.T) {
		t.Run("user2 is friend of user1", func(t *testing.T) {
			url := testServer.URL + "/friends/" + strconv.Itoa(user1.ID)

			res, err := http.Get(url)
			chat.AssertTestErr(t, err)

			resBody, err := ioutil.ReadAll(res.Body)
			chat.AssertTestErr(t, err)

			chat.AssertTestStatusCode(t, res.StatusCode, http.StatusOK)

			var decodedResBody []map[string]interface{}
			err = json.Unmarshal(resBody, &decodedResBody)
			chat.AssertTestErr(t, err)

			want := []users.User{user2}
			got := decodedResBody

			if !reflect.DeepEqual(got, want) {
				t.Errorf("wanted %#v, got %#v", want, got)
			}
		})

		t.Run("user1 is friend of user2", func(t *testing.T) {
			url := testServer.URL + "/friends/" + strconv.Itoa(user2.ID)

			res, err := http.Get(url)
			chat.AssertTestErr(t, err)

			resBody, err := ioutil.ReadAll(res.Body)
			chat.AssertTestErr(t, err)

			chat.AssertTestStatusCode(t, res.StatusCode, http.StatusOK)

			var decodedResBody []map[string]interface{}
			err = json.Unmarshal(resBody, &decodedResBody)
			chat.AssertTestErr(t, err)

			want := []users.User{user1}
			got := decodedResBody

			if !reflect.DeepEqual(got, want) {
				t.Errorf("wanted %#v, got %#v", want, got)
			}
		})
	})
}
