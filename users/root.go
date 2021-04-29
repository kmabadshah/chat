package users

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/kmabadshah/chat"
	"log"
	"testing"
)

type User struct {
	Uname string `json:"uname"`
	Pass  string `json:"pass"`
	ID    int    `json:"id"`
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
	router.Path("/users").Methods("POST").HandlerFunc(AddUser)
	router.Path("/users/{id}").Methods("GET").HandlerFunc(GetUser)
	router.Path("/users").Methods("GET").HandlerFunc(GetUsers)
	router.Path("/users/{id}").Methods("PUT").HandlerFunc(UpdateUser)

	return router
}

func CreateTestUser(t *testing.T) User {
	user := User{
		Uname: chat.RandSeq(5),
		Pass:  chat.RandSeq(5),
	}

	_, err := chat.DB.Model(&user).Insert()
	chat.AssertTestErr(t, err)

	return user
}
