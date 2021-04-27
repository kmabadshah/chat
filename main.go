package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
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
	router.Path("/users").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := ioutil.ReadAll(r.Body)
		if assertRandomError(err, &w) {
			return
		}

		var decodedReqBody map[string]string
		err = json.Unmarshal(reqBody, &decodedReqBody)
		if err != nil || decodedReqBody["uname"] == "" || decodedReqBody["pass"] == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user := User{
			Uname: decodedReqBody["uname"],
			Pass:  decodedReqBody["pass"],
		}
		res, err := db.Model(&user).Insert()
		if err != nil || res.RowsAffected() != 1 || user.ID == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		encodedResBody, err := json.Marshal(user)
		if assertRandomError(err, &w) {
			return
		}

		_, err = w.Write(encodedResBody)
		if assertRandomError(err, &w) {
			return
		}
	})

	return router
}

func assertError(err error, w *http.ResponseWriter, statusCode int) bool {
	if err != nil {
		(*w).WriteHeader(statusCode)
		fmt.Println(err)
		return true
	}

	return false
}

func assertRandomError(err error, w *http.ResponseWriter) bool {
	return assertError(err, w, http.StatusInternalServerError)
}

func init() {
	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		panic(err)
	}
}
