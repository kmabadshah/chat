package messages

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/kmabadshah/chat"
	"log"
)

type Friend struct {
	SrcID int `json:"srcID"`
	TarID int `json:"tarID"`
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

	return router
}
