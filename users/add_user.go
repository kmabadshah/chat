package users

import (
	"encoding/json"
	"fmt"
	"github.com/kmabadshah/chat/shared"
	"io/ioutil"
	"log"
	"net/http"
)

func AddUser(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)

	if shared.AssertInternalError(err, &w) {
		return
	}

	var decodedReqBody map[string]string
	err = json.Unmarshal(reqBody, &decodedReqBody)

	if err != nil || decodedReqBody["uname"] == "" || decodedReqBody["pass"] == "" {
		log.Println(err, decodedReqBody)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if already exists
	var usr []User
	err = shared.DB.Model(&usr).Where("uname=?", decodedReqBody["uname"]).Select()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(usr) == 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := User{
		Uname: decodedReqBody["uname"],
		Pass:  decodedReqBody["pass"],
	}
	res, err := shared.DB.Model(&user).Insert()
	if err != nil || res.RowsAffected() != 1 || user.ID == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
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
