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
		sendGetFriendsReqAndCompareRes := func(id int, wantedResBody []map[string]int) {
			url := testServer.URL + "/friends/" + strconv.Itoa(id)

			res, err := http.Get(url)
			chat.AssertTestErr(t, err)

			resBody, err := ioutil.ReadAll(res.Body)
			chat.AssertTestErr(t, err)

			chat.AssertTestStatusCode(t, res.StatusCode, http.StatusOK)

			var decodedResBody []map[string]int
			err = json.Unmarshal(resBody, &decodedResBody)
			chat.AssertTestErr(t, err)

			if !reflect.DeepEqual(wantedResBody, decodedResBody) {
				t.Errorf("wanted %#v, got %#v", wantedResBody, decodedResBody)
			}
		}

		t.Run("user2 is friend of user1", func(t *testing.T) {
			resBody := []map[string]int{
				{"srcID": user1.ID, "tarID": user2.ID},
			}
			sendGetFriendsReqAndCompareRes(user1.ID, resBody)
		})

		t.Run("user1 is friend of user2", func(t *testing.T) {
			resBody := []map[string]int{
				{"srcID": user1.ID, "tarID": user2.ID},
			}
			sendGetFriendsReqAndCompareRes(user2.ID, resBody)
		})
	})
}
