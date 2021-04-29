package users

import (
	"encoding/json"
	"github.com/kmabadshah/chat"
	"net/http"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	err := chat.DB.Model(&users).Select()
	chat.AssertInternalError(err, &w)

	encodedResBody, err := json.Marshal(users)
	if chat.AssertInternalError(err, &w) {
		return
	}

	_, err = w.Write(encodedResBody)
	if chat.AssertInternalError(err, &w) {
		return
	}
}
