package messages

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

func TestGetMessages(t *testing.T) {
	shared.ClearAllTables()

	router := NewRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	user1 := users.CreateTestUser(t)
	user2 := users.CreateTestUser(t)

	// add message
	message := Message{
		SrcID: user1.ID,
		TarID: user2.ID,
		Text:  "hello user2",
	}
	res, err := shared.DB.Model(&message).Insert()
	if err != nil || res.RowsAffected() != 1 {
		t.Fatal(err)
	}

	t.Run("valid id", func(t *testing.T) {
		url := testServer.URL + "/messages/" + strconv.Itoa(user1.ID)
		res, err := http.Get(url)
		shared.AssertTestErr(t, err)

		resBody, err := ioutil.ReadAll(res.Body)
		shared.AssertTestErr(t, err)

		shared.AssertTestStatusCode(t, res.StatusCode, http.StatusOK)

		var decodedResBody []map[string]interface{}
		err = json.Unmarshal(resBody, &decodedResBody)
		shared.AssertTestErr(t, err)

		want := []map[string]interface{}{
			{
				"srcID": float64(user1.ID),
				"tarID": float64(user2.ID),
				"text":  "hello user2",
			},
		}
		got := decodedResBody

		if !reflect.DeepEqual(got, want) {
			t.Errorf("want %#v, got %#v", want, got)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		url := testServer.URL + "/messages/" + strconv.Itoa(-1)
		res, err := http.Get(url)
		shared.AssertTestErr(t, err)

		shared.AssertTestStatusCode(t, res.StatusCode, http.StatusNotFound)
	})
}
