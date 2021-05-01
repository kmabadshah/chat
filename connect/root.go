package connect

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/kmabadshah/chat/shared"
	"github.com/kmabadshah/chat/users"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile)

	ctx := context.Background()
	if err := shared.DB.Ping(ctx); err != nil {
		panic(err)
	}
}

type Client struct {
	conn *websocket.Conn
	uid  int
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	router.Path("/connect/{uid}").Methods("GET").HandlerFunc(HandleConnect)

	return router
}

func testCreateAndConnect(t *testing.T, url string) (users.User, *websocket.Conn) {
	t.Helper()

	user := users.CreateTestUser(t)
	url1 := "ws" + strings.TrimPrefix(url, "http") + "/connect/" + strconv.Itoa(user.ID)
	conn, res, err := websocket.DefaultDialer.Dial(url1, nil)
	shared.AssertTestErr(t, err)
	shared.AssertTestStatusCode(t, res.StatusCode, http.StatusSwitchingProtocols)

	return user, conn
}

func testSendMsg(t *testing.T, reqBody map[string]interface{}, conn *websocket.Conn) {
	t.Helper()

	encodedReqBody, err := json.Marshal(reqBody)
	shared.AssertTestErr(t, err)

	err = conn.WriteMessage(websocket.TextMessage, encodedReqBody)
	shared.AssertTestErr(t, err)
}

func testRecvMsg(t *testing.T, conn *websocket.Conn) map[string]interface{} {
	t.Helper()

	_, encodedData, err := conn.ReadMessage()
	var decodedData map[string]interface{}
	err = json.Unmarshal(encodedData, &decodedData)
	shared.AssertTestErr(t, err)

	return decodedData
}

func testReadMsgAndComp(t *testing.T, conn *websocket.Conn, want map[string]interface{}) {
	t.Helper()

	decData := testRecvMsg(t, conn)

	if !reflect.DeepEqual(decData, want) {
		t.Errorf("got %#v, wanted %#v", decData, want)
	}
}
