package users

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kmabadshah/chat/shared"
	"net/http"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	uname := mux.Vars(r)["uname"]

	var user User
	err := shared.DB.Model(&user).Where("uname=?", uname).Select()
	if err != nil || user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	encodedResBody, err := json.Marshal(user)
	if shared.AssertInternalError(err, &w) {
		return
	}

	_, err = w.Write(encodedResBody)
	if shared.AssertInternalError(err, &w) {
		return
	}
}
