package friends

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kmabadshah/chat"
	"log"
	"net/http"
)

func GetFriends(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var friends []Friend
	err := chat.DB.Model(&friends).Where("src_id=?", id).WhereOr("tar_id=?", id).Select()
	if err != nil {
		log.Println(err)
	}

	encodedResBody, err := json.Marshal(friends)
	if chat.AssertRandomError(err, &w) {
		return
	}

	_, err = w.Write(encodedResBody)
	if err != nil {
		log.Println(err)
	}
}
