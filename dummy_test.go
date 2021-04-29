package chat

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/kmabadshah/chat/users"
	"testing"
)

func TestDummy(t *testing.T) {
	x := map[string]interface{}{
		"uname":     "adnan",
		"pass":      "hello",
		"id":        "helo",
		"something": "else",
	}
	encodedX, err := json.Marshal(x)
	AssertTestErr(t, err)

	var user users.User
	err = json.Unmarshal(encodedX, &user)
	AssertTestErr(t, err)

	//var y users.User
	//err = mapstructure.Decode(user, &y)
	//AssertTestErr(t, err)

	spew.Dump(user)
}
