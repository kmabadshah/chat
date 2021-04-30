package friends

import (
	"encoding/json"
	"github.com/kmabadshah/chat"
	"io/ioutil"
	"net/http"
)

const errReqBody = "invalid body, must include a valid srcID and tarID field"

func AddFriend(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if chat.AssertInternalError(err, &w) {
		return
	}

	var decodedReqBody map[string]int
	err = json.Unmarshal(reqBody, &decodedReqBody)
	if chat.AssertInternalError(err, &w) {
		return
	}

	if len(decodedReqBody) != 2 ||
		decodedReqBody["srcID"] == 0 ||
		decodedReqBody["tarID"] == 0 {

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(errReqBody))

		return
	}

	friend := Friend{
		SrcID: decodedReqBody["srcID"],
		TarID: decodedReqBody["tarID"],
	}
	res, err := chat.DB.Model(&friend).Insert()
	if err != nil || res.RowsAffected() != 1 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(errReqBody))
		return
	}
}
