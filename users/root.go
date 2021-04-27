package users

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
	"net/http"
	"testing"
)

var db = pg.Connect(&pg.Options{
	Addr:     ":5432",
	User:     "kmab",
	Password: "kmab",
	Database: "chat_test",
})

type User struct {
	Uname string `json:"uname"`
	Pass  string `json:"pass"`
	ID    int    `json:"id"`
}

func newRouter() *mux.Router {
	router := mux.NewRouter()
	router.Path("/users").Methods("POST").HandlerFunc(addUser)

	return router
}

func init() {
	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		panic(err)
	}
}

func assertError(err error, w *http.ResponseWriter, statusCode int) bool {
	if err != nil {
		(*w).WriteHeader(statusCode)
		fmt.Println(err)
		return true
	}

	return false
}

func assertTestErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func assertRandomError(err error, w *http.ResponseWriter) bool {
	return assertError(err, w, http.StatusInternalServerError)
}
