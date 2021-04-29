package users

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kmabadshah/chat"
	"net/http"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	uname := mux.Vars(r)["uname"]

	var user User
	err := chat.DB.Model(&user).Where("uname=?", uname).Select()
	if err != nil || user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	encodedResBody, err := json.Marshal(user)
	if chat.AssertInternalError(err, &w) {
		return
	}

	_, err = w.Write(encodedResBody)
	if chat.AssertInternalError(err, &w) {
		return
	}
}
