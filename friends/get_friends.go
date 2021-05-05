package friends

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kmabadshah/chat/shared"
	"github.com/kmabadshah/chat/users"
	"log"
	"net/http"
)

func GetFriends(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// check if valid id
	var users []users.User
	shared.DB.Model(&users).Where("id=?", id).Select()

	if len(users) != 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var friends []Friend
	err := shared.DB.Model(&friends).Where("src_id=?", id).WhereOr("tar_id=?", id).Select()
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//if len(friends) == 0 {
	//	w.WriteHeader(http.StatusNotFound)
	//	return
	//}

	encodedResBody, err := json.Marshal(friends)
	if shared.AssertInternalError(err, &w) {
		return
	}

	_, err = w.Write(encodedResBody)
	if err != nil {
		log.Println(err)
	}
}
