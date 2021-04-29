package friends

import (
	"encoding/json"
	"github.com/kmabadshah/chat"
	"io/ioutil"
	"log"
	"net/http"
)

func AddFriend(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if chat.AssertRandomError(err, &w) {
		return
	}

	var decodedReqBody map[string]int
	err = json.Unmarshal(reqBody, &decodedReqBody)
	if chat.AssertRandomError(err, &w) {
		return
	}

	friend := Friend{
		SrcID: decodedReqBody["srcID"],
		TarID: decodedReqBody["tarID"],
	}
	res, err := chat.DB.Model(&friend).Insert()
	if err != nil || res.RowsAffected() != 1 {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
