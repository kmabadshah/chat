package messages

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/kmabadshah/chat"
	"log"
)

type Message struct {
	SrcID int    `json:"srcID" mapstructure:"srcID"`
	TarID int    `json:"tarID" mapstructure:"tarID"`
	Text  string `json:"text" mapstructure:"text"`
}

func init() {
	log.SetFlags(log.Lshortfile)

	ctx := context.Background()
	if err := chat.DB.Ping(ctx); err != nil {
		panic(err)
	}
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.Path("/messages").Methods("POST").HandlerFunc(AddMessage)

	return router
}
