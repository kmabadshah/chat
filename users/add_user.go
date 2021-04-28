package users

import (
	"encoding/json"
	"fmt"
	"github.com/kmabadshah/chat"
	"io/ioutil"
	"net/http"
)

func AddUser(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)

	if chat.AssertRandomError(err, &w) {
		return
	}

	var decodedReqBody map[string]string
	err = json.Unmarshal(reqBody, &decodedReqBody)
	if err != nil || decodedReqBody["uname"] == "" || decodedReqBody["pass"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := User{
		Uname: decodedReqBody["uname"],
		Pass:  decodedReqBody["pass"],
	}
	res, err := chat.DB.Model(&user).Insert()
	if err != nil || res.RowsAffected() != 1 || user.ID == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
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
