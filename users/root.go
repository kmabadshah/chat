package users

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/kmabadshah/chat/shared"
	"log"
	"testing"
)

type User struct {
	Uname string `json:"uname" mapstructure:"uname"`
	Pass  string `json:"pass" mapstructure:"pass"`
	ID    int    `json:"id" mapstructure:"id"`
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
	router.Path("/users").Methods("POST").HandlerFunc(AddUser)
	router.Path("/users/{uname}").Methods("GET").HandlerFunc(GetUser)
	router.Path("/users").Methods("GET").HandlerFunc(GetUsers)
	router.Path("/users/{id}").Methods("PUT").HandlerFunc(UpdateUser)

	return router
}

func CreateTestUser(t *testing.T) User {
	user := User{
		Uname: shared.RandSeq(5),
		Pass:  shared.RandSeq(5),
	}

	_, err := shared.DB.Model(&user).Insert()
	shared.AssertTestErr(t, err)

	return user
}
