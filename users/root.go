package users

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"testing"
	"time"
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
	router.Path("/users/{id}").Methods("GET").HandlerFunc(getUser)
	router.Path("/users").Methods("GET").HandlerFunc(getUsers)
	router.Path("/users/{id}").Methods("PUT").HandlerFunc(updateUser)

	return router
}

func init() {
	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		panic(err)
	}
}

func clearUserTable() {
	_, err := db.Exec(`delete from users`)
	if err != nil {
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

func randSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func createTestUser(t *testing.T) User {
	user := User{
		Uname: randSeq(5),
		Pass:  randSeq(5),
	}

	_, err := db.Model(&user).Insert()
	assertTestErr(t, err)

	return user
}

func assertTestErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func assertTestStatusCode(t *testing.T, got int, want int) {
	t.Helper()

	if got != want {
		t.Errorf("invalid code, wanted %#v, got %#v", want, got)
	}
}

func assertRandomError(err error, w *http.ResponseWriter) bool {
	return assertError(err, w, http.StatusInternalServerError)
}
