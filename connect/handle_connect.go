package connect

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/kmabadshah/chat"
	"log"
	"net/http"
	"strconv"
)

var clients []Client

func HandleConnect(w http.ResponseWriter, r *http.Request) {
	uidRaw := mux.Vars(r)["uid"]
	uid, err := strconv.Atoi(uidRaw)
	chat.AssertInternalError(err, &w)

	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	chat.AssertInternalError(err, &w)

	defer conn.Close()

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
			for _, cl := range clients {
				encMsg, _ := json.Marshal(map[string]interface{}{
					"type": "leave",
					"id":   uid,
				})
				err = cl.conn.WriteMessage(websocket.TextMessage, encMsg)
				if err != nil {
					log.Println(err)
				}
			}
			break
		}
		err = json.Unmarshal(encodedData, &decodedData)
		if err != nil {
			log.Println(err)
			break
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
}
