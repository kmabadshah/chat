package chat

import (
	"encoding/json"
	"fmt"
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

type X struct {
	clients []Client
	router  http.Handler
}

func init() {
	log.SetFlags(0)
}

func newX() *X {
	x := X{}
	upgrader := websocket.Upgrader{}
	router := mux.NewRouter()
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
		x.clients = append(x.clients, clMain)

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
				for _, client := range x.clients {
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
				for i, client := range x.clients {
					if reflect.DeepEqual(client.conn, conn) {
						x.clients = append(x.clients[:i], x.clients[i+1:]...)
					}
				}

				//prettyError(err)
				break
			}
		}

	})

	x.router = router

	return &x
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
