package users

import (
	"encoding/json"
	"github.com/kmabadshah/chat/shared"
	"net/http"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	err := shared.DB.Model(&users).Select()
	shared.AssertInternalError(err, &w)

	encodedResBody, err := json.Marshal(users)
	if shared.AssertInternalError(err, &w) {
		return
	}

	_, err = w.Write(encodedResBody)
	if shared.AssertInternalError(err, &w) {
		return
	}
}
