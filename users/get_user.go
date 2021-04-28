package users

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kmabadshah/chat"
	"net/http"
	"strconv"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	idRaw := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idRaw)
	if chat.AssertError(err, &w, http.StatusNotFound) {
		return
	}

	var user User
	err = chat.DB.Model(&user).Where("id=?", id).Select()
	if err != nil || user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	encodedResBody, err := json.Marshal(user)
	if chat.AssertRandomError(err, &w) {
		return
	}

	_, err = w.Write(encodedResBody)
	if chat.AssertRandomError(err, &w) {
		return
	}
}
