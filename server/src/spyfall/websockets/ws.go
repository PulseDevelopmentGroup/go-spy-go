package websockets

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

var ClientById = make(map[string]*websocket.Conn)
var ClientByConn = make(map[*websocket.Conn]string)

//Websockets Upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Request struct {
	Kind string `json:"kind"`
	Data string `json:"data"`
}

type Response struct {
	Kind string `json:"kind"`
	Data string `json:"data"`
	Err  string `json:"error,omitempty"`
}

type GameData struct {
	GameID   string `json:"gameid"`
	Username string `json:"username"`
}

type LeaveData struct {
	Username string `json:"username"`
	Reason   string `json:"reason,omitempty"`
}

type ErrData struct {
	Err  string `json:"error"`
	Desc string `json:"description,omitempty"`
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func SendToPlayer(response *Response, connection *websocket.Conn) (error, error) {
	r, marshalErr := json.Marshal(response)
	writeErr := connection.WriteMessage(1, r)
	return marshalErr, writeErr
}
