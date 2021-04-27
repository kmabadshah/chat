package users

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func getUser(w http.ResponseWriter, r *http.Request) {
	idRaw := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idRaw)
	if assertError(err, &w, http.StatusNotFound) {
		return
	}

	var user User
	err = db.Model(&user).Where("id=?", id).Select()
	if err != nil || user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	encodedResBody, err := json.Marshal(user)
	if assertRandomError(err, &w) {
		return
	}

	_, err = w.Write(encodedResBody)
	if assertRandomError(err, &w) {
		return
	}
}
