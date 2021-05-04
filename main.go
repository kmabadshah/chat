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

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			next.ServeHTTP(w, r)
		})
	})

	router.PathPrefix("/api/").Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Path("/api/connect/{uid}").Methods("GET").HandlerFunc(connect.HandleConnect)

	router.Path("/api/friends").Methods("POST").HandlerFunc(friends.AddFriend)
	router.Path("/api/friends/{id}").Methods("GET").HandlerFunc(friends.GetFriends)

	router.Path("/api/messages").Methods("POST").HandlerFunc(messages.AddMessage)
	router.Path("/api/messages/{uid}").Methods("GET").HandlerFunc(messages.GetMessages)

	router.Path("/api/users").Methods("POST").HandlerFunc(users.AddUser)
	router.Path("/api/users/{uname}").Methods("GET").HandlerFunc(users.GetUser)
	router.Path("/api/users").Methods("GET").HandlerFunc(users.GetUsers)
	router.Path("/api/users/{id}").Methods("PUT").HandlerFunc(users.UpdateUser)

	fmt.Println("listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}
