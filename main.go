package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"runtime"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	name string
}

type Hub struct {
	clients []Client
	router  http.Handler
}

type User struct {
	ID    uint   `json:"id"`
	Uname string `json:"uname"`
	Pass  string `json:"pass"`
}

var db = pg.Connect(&pg.Options{
	Addr:     ":5432",
	User:     "kmab",
	Password: "kmab",
	Database: "chat_test",
})

func init() {
	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		panic(err)
	}
}

func newHub() *Hub {
	hub := Hub{}
	upgrader := websocket.Upgrader{}
	router := mux.NewRouter()

	router.PathPrefix("/users/{uname}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the user with {name}
		uname := mux.Vars(r)["uname"]
		// query the db
		var user User
		err := db.Model(&user).Where("uname=?", uname).Select()
		if assertError(err) {
			return
		}

		if user.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("invalid username"))
			return
		}

		// encode and send user
		encodedResBody, err := json.Marshal(user)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(encodedResBody)
	})
	router.PathPrefix("/users").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the data from req body and decode
		reqBody, err := ioutil.ReadAll(r.Body)
		var decodedReqBody map[string]interface{}
		err = json.Unmarshal(reqBody, &decodedReqBody)
		isValidBody := decodedReqBody["uname"] != "" || decodedReqBody["pass"] != ""
		if err != nil || !isValidBody {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("invalid req body, must contain a uname and pass field"))
			return
		}

		// create the user
		user := User{
			Uname: decodedReqBody["uname"].(string),
			Pass:  decodedReqBody["pass"].(string),
		}
		res, err := db.Model(&user).Insert()
		if err != nil || res.RowsAffected() != 1 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("insertErr: ", err)
			return
		}

		// encode and send the user
		resBody, _ := json.Marshal(user)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resBody)
	})
	router.PathPrefix("/connect/{name}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]

		// upgrade the req
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			prettyError(err)
			return
		}
		defer conn.Close()

		// add the client
		clMain := Client{conn: conn, name: name}
		hub.clients = append(hub.clients, clMain)

		// notify the client
		err = conn.WriteMessage(websocket.TextMessage, []byte("connection successful"))
		if assertError(err) {
			return
		}

		for {
			_, msg, err := conn.ReadMessage()

			if string(msg) != "" {
				// decode
				var decodedMsg map[string]string
				err = json.Unmarshal(msg, &decodedMsg)
				if err != nil {
					prettyError(err)
					return
				}

				// when msg recv from client1 send it to client2
				receiver := decodedMsg["receiver"]
				for _, client := range hub.clients {
					if client.name == receiver {
						err = client.conn.WriteMessage(websocket.TextMessage, []byte(decodedMsg["message"]))
						if assertError(err) {
							return
						}
					}
				}
			}

			// remove the client on disconnect
			if err != nil {
				for i, client := range hub.clients {
					if reflect.DeepEqual(client.conn, conn) {
						hub.clients = append(hub.clients[:i], hub.clients[i+1:]...)
					}
				}

				//prettyError(err)
				break
			}
		}

	})

	hub.router = router

	return &hub
}

func prettyLog(v ...interface{}) {
	_, fileName, line, _ := runtime.Caller(1)
	fmt.Printf("[log] %s:%d\n", fileName, line)
	spew.Dump(v...)
}

func prettyError(v ...interface{}) {
	_, fileName, line, _ := runtime.Caller(1)
	fmt.Printf("[error] %s:%d ", fileName, line)
	fmt.Println(v...)
}

func assertError(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}

	return false
}

func main() {
	db := pg.Connect(&pg.Options{
		Addr:     ":5432",
		User:     "kmab",
		Password: "kmab",
		Database: "chat_test",
	})
	defer db.Close()

	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		panic(err)
	}

	// delete all users
	_, err := db.Model((*User)(nil)).Exec(`
	    delete from users
	`)
	if err != nil {
		panic(err)
	}

	// insert a user
	user1 := &User{
		Uname: "adnan",
		Pass:  "badshah",
	}
	_, err = db.Model(user1).Insert()
	if err != nil {
		panic(err)
	}

	// get the user
	var user User
	db.Model(&user).Where(`id=?`, user1.ID).Select()

	fmt.Println(user)
}

func init() {
	log.SetFlags(0)
}
