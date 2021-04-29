package messages

import (
	"encoding/json"
	"github.com/kmabadshah/chat"
	"io/ioutil"
	"net/http"
)

const resAddSuccess = "message added successfully"

func AddMessage(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var decodedReqBody map[string]interface{}
	_ = json.Unmarshal(reqBody, &decodedReqBody)

	message := Message{
		SrcID: int(decodedReqBody["srcID"].(float64)),
		TarID: int(decodedReqBody["tarID"].(float64)),
		Text:  decodedReqBody["text"].(string),
	}
	_, _ = chat.DB.Model(&message).Insert()

	_, _ = w.Write([]byte(resAddSuccess))
}
