package connect

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/kmabadshah/chat"
	"log"
	"net/http"
	"strconv"
)

func init() {
	log.SetFlags(log.Lshortfile)

	ctx := context.Background()
	if err := chat.DB.Ping(ctx); err != nil {
		panic(err)
	}
}

type Client struct {
	conn *websocket.Conn
	uid  int
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	var clients []Client

	router.Path("/connect/{uid}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uidRaw := mux.Vars(r)["uid"]
		uid, err := strconv.Atoi(uidRaw)
		chat.AssertInternalError(err, &w)

		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		chat.AssertInternalError(err, &w)

		for _, cl := range clients {
			encodedData, err := json.Marshal(map[string]interface{}{
				"type": "connect",
				"id":   uid,
			})
			err = cl.conn.WriteMessage(websocket.TextMessage, encodedData)
			if err != nil {
				log.Println(err)
			}
		}

		clients = append(clients, Client{
			conn, uid,
		})

		for {
			var decodedData map[string]interface{}
			_, encodedData, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
			}
			err = json.Unmarshal(encodedData, &decodedData)
			if err != nil {
				log.Println(err)
			}

			switch decodedData["type"].(string) {
			case "message":
				tarIDRaw := decodedData["uid"].(float64)
				tarID := int(tarIDRaw)

				encodedData, _ = json.Marshal(map[string]interface{}{
					"type": "message",
					"uid":  uid,
				})

				for _, cl := range clients {
					if cl.uid == tarID {
						_ = cl.conn.WriteMessage(websocket.TextMessage, encodedData)
					}
				}

			case "broadcast":
				encodedData, _ = json.Marshal(map[string]interface{}{
					"type": "user",
					"uid":  uid,
				})

				for _, cl := range clients {
					err = cl.conn.WriteMessage(websocket.TextMessage, encodedData)
					if err != nil {
						log.Println(err)
					}
				}
			}

		}
	})

	return router
}
