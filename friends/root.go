package friends

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/kmabadshah/chat/shared"
	"log"
)

type Friend struct {
	SrcID int `json:"srcID"`
	TarID int `json:"tarID"`
}

func init() {
	log.SetFlags(log.Lshortfile)

	ctx := context.Background()
	if err := shared.DB.Ping(ctx); err != nil {
		panic(err)
	}
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.Path("/friends").Methods("POST").HandlerFunc(AddFriend)
	router.Path("/friends/{id}").Methods("GET").HandlerFunc(GetFriends)

	return router
}
