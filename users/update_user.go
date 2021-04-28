package users

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

var errUnameInUse = "that username is already in use"

func updateUser(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid body, must contain a valid uname and/or pass field"))
		return
	}

	var decodedReqBody map[string]interface{}
	err = json.Unmarshal(reqBody, &decodedReqBody)
	if err != nil || (decodedReqBody["pass"] == nil && decodedReqBody["uname"] == nil) {
		log.Println(err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid body, must contain a valid uname and/or pass field"))
		return
	}

	// check if the uname exists
	var userCheck User
	err = db.Model(&userCheck).Where("uname=?", decodedReqBody["uname"]).Select()
	if userCheck.ID != 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errUnameInUse))
		return
	}

	// update the user
	var user User
	id := mux.Vars(r)["id"]
	res, err := db.Model(&decodedReqBody).TableExpr("users").Where("id=?", id).Update()
	if err != nil || res.RowsAffected() != 1 {
		log.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.Model(&user).Where("id=?", id).Select()
	if err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
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
