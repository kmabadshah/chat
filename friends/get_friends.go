package friends

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kmabadshah/chat/shared"
	"github.com/kmabadshah/chat/users"
	"log"
	"net/http"
	"strconv"
)

func GetFriends(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// check if valid id
	var usrs []users.User
	shared.DB.Model(&usrs).Where("id=?", id).Select()

	if len(usrs) != 1 {
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

	var friendUsers []users.User

	for _, f := range friends {
		intId, _ := strconv.Atoi(id)

		var user users.User
		var err error

		if f.SrcID != intId {
			err = shared.DB.Model(&user).Where("id=?", f.SrcID).Select()
		} else {
			err = shared.DB.Model(&user).Where("id=?", f.TarID).Select()
		}

		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		friendUsers = append(friendUsers, user)
	}

	encodedResBody, err := json.Marshal(friendUsers)
	if shared.AssertInternalError(err, &w) {
		return
	}

	_, err = w.Write(encodedResBody)
	if err != nil {
		log.Println(err)
	}
}
