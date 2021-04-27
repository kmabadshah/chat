package users

import (
	"encoding/json"
	"net/http"
)

func getUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	err := db.Model(&users).Select()
	assertRandomError(err, &w)

	encodedResBody, err := json.Marshal(users)
	if assertRandomError(err, &w) {
		return
	}

	_, err = w.Write(encodedResBody)
	if assertRandomError(err, &w) {
		return
	}
}
