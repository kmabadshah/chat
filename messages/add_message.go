package messages

import (
	"encoding/json"
	"github.com/kmabadshah/chat"
	"io/ioutil"
	"net/http"
)

const (
	resAddSuccess = "message added successfully"
	errReqBody    = "invalid req body, must include valid srcID, tarID and text field"
)

func AddMessage(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if chat.AssertInternalError(err, &w) {
		return
	}

	var decodedReqBody map[string]interface{}
	err = json.Unmarshal(reqBody, &decodedReqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(errReqBody))
		return
	}

	if len(decodedReqBody) != 3 ||
		decodedReqBody["srcID"] == nil ||
		decodedReqBody["tarID"] == nil ||
		decodedReqBody["text"] == nil {

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(errReqBody))
		return
	}

	message := Message{
		SrcID: int(decodedReqBody["srcID"].(float64)),
		TarID: int(decodedReqBody["tarID"].(float64)),
		Text:  decodedReqBody["text"].(string),
	}
	_, _ = chat.DB.Model(&message).Insert()

	_, _ = w.Write([]byte(resAddSuccess))
}
