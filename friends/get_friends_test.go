package friends

import (
	"encoding/json"
	"github.com/kmabadshah/chat/shared"
	"github.com/kmabadshah/chat/users"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

func TestGetFriends(t *testing.T) {
	shared.ClearAllTables()

	router := NewRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	user1 := users.CreateTestUser(t)
	user2 := users.CreateTestUser(t)

	sendReqAndCompareRes := func(t *testing.T, id int, wantedStatus int, wantedResBody interface{}) {
		t.Helper()
		url := testServer.URL + "/friends/" + strconv.Itoa(id)

		res, err := http.Get(url)
		shared.AssertTestErr(t, err)

		resBody, err := ioutil.ReadAll(res.Body)
		shared.AssertTestErr(t, err)

		shared.AssertTestStatusCode(t, res.StatusCode, wantedStatus)

		if wantedResBody != nil {
			var decodedResBody []map[string]interface{}
			err = json.Unmarshal(resBody, &decodedResBody)
			shared.AssertTestErr(t, err)

			if !reflect.DeepEqual(wantedResBody, decodedResBody) {
				t.Errorf("wanted %#v, got %#v", wantedResBody, decodedResBody)
			}
		}
	}

	t.Run("valid uid with out friends", func(t *testing.T) {
		sendReqAndCompareRes(t, user1.ID, 200, nil)
		sendReqAndCompareRes(t, user2.ID, 200, nil)
	})

	t.Run("valid uid with existing friend", func(t *testing.T) {
		// add a friend
		friend := Friend{SrcID: user1.ID, TarID: user2.ID}
		result, err := shared.DB.Model(&friend).Insert()
		if err != nil || result.RowsAffected() != 1 {
			t.Fatal(err)
		}

		t.Run("user2 is friend of user1", func(t *testing.T) {
			resBody := []map[string]interface{}{
				{"uname": user2.Uname, "pass": user2.Pass, "id": float64(user2.ID)},
			}
			sendReqAndCompareRes(t, user1.ID, 200, resBody)
		})

		t.Run("user1 is friend of user2", func(t *testing.T) {
			resBody := []map[string]interface{}{
				{"uname": user1.Uname, "pass": user1.Pass, "id": float64(user1.ID)},
			}
			sendReqAndCompareRes(t, user2.ID, 200, resBody)
		})
	})

	t.Run("invalid id", func(t *testing.T) {
		sendReqAndCompareRes(t, -1, http.StatusNotFound, nil)
	})
}
