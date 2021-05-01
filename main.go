package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kmabadshah/chat/connect"
	"github.com/kmabadshah/chat/friends"
	"github.com/kmabadshah/chat/messages"
	"github.com/kmabadshah/chat/shared"
	"github.com/kmabadshah/chat/users"
	"log"
	"net/http"
)

func main() {
	shared.ClearAllTables()

	router := mux.NewRouter()

	router.Path("/connect/{uid}").Methods("GET").HandlerFunc(connect.HandleConnect)

	router.Path("/friends").Methods("POST").HandlerFunc(friends.AddFriend)
	router.Path("/friends/{id}").Methods("GET").HandlerFunc(friends.GetFriends)

	router.Path("/messages").Methods("POST").HandlerFunc(messages.AddMessage)
	router.Path("/messages/{uid}").Methods("GET").HandlerFunc(messages.GetMessages)

	router.Path("/users").Methods("POST").HandlerFunc(users.AddUser)
	router.Path("/users/{uname}").Methods("GET").HandlerFunc(users.GetUser)
	router.Path("/users").Methods("GET").HandlerFunc(users.GetUsers)
	router.Path("/users/{id}").Methods("PUT").HandlerFunc(users.UpdateUser)

	fmt.Println("listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
