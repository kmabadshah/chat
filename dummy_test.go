package main

import (
	"encoding/json"
	"github.com/kmabadshah/chat/messages"
	"github.com/kmabadshah/chat/shared"
	"testing"
)

type y struct {
	names string
}

func TestDummy(t *testing.T) {
	x := map[string]interface{}{
		"srcID": 10,
		"tarID": 20,
		"text":  "hello world",
	}
	encodedX, err := json.Marshal(x)
	shared.AssertTestErr(t, err)

	var decodedX messages.Message
	err = json.Unmarshal(encodedX, &decodedX)
	shared.AssertTestErr(t, err)

}
