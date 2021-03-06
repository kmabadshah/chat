package shared

import (
	"github.com/go-pg/pg/v10"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
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

	_, err = DB.Exec(`delete from messages`)
	if err != nil {
		panic(err)
	}
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

func AssertInternalError(err error, w *http.ResponseWriter) bool {
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)

		_, fullPath, linum, _ := runtime.Caller(1)
		fileName := filepath.Base(fullPath)
		customLogger := log.New(os.Stderr, fileName+":"+strconv.Itoa(linum)+" ", 0)
		customLogger.Println(err)

		return true
	}

	return false
}
