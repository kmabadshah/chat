package chat

import (
	"github.com/go-pg/pg/v10"
	"log"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

var DB = pg.Connect(&pg.Options{
	Addr:     ":5432",
	User:     "kmab",
	Password: "kmab",
	Database: "chat_test",
})

func ClearAllTables() {
	_, err := DB.Exec(`delete from users`)
	if err != nil {
		panic(err)
	}

	_, err = DB.Exec(`delete from friends`)
	if err != nil {
		panic(err)
	}
}

func AssertError(err error, w *http.ResponseWriter, statusCode int) bool {
	if err != nil {
		(*w).WriteHeader(statusCode)
		log.Println(err)
		return true
	}

	return false
}

func RandSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func AssertTestErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func AssertTestStatusCode(t *testing.T, got int, want int) {
	t.Helper()

	if got != want {
		t.Errorf("invalid code, wanted %#v, got %#v", want, got)
	}
}

func AssertRandomError(err error, w *http.ResponseWriter) bool {
	return AssertError(err, w, http.StatusInternalServerError)
}
