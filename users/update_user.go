package users

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"net/http"
)

func updateUser(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid body, must contain a valid uname and/or pass field"))
		return
	}

	var decodedReqBody map[string]string
	err = json.Unmarshal(reqBody, &decodedReqBody)
	if err != nil || decodedReqBody["pass"] == "" {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid body, must contain a valid uname and/or pass field"))
		return
	}

	var user User
	id := mux.Vars(r)["id"]
	err = mapstructure.Decode(decodedReqBody, &user)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid body, must contain a valid uname and/or pass field"))
		return
	}

	res, err := db.Model(&user).Column("pass").Where("id=?", id).Update()
	if err != nil || res.RowsAffected() != 1 {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.Model(&user).Where("id=?", id).Select()
	if err != nil {
		fmt.Println(err)
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
