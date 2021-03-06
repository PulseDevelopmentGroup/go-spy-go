package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

//Websockets Upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type request struct {
	Kind string `json:"kind"`
	Data string `json:"data"`
}

type response struct {
	Kind string `json:"kind"`
	Data string `json:"data"`
	Err  string `json:"error,omitempty"`
}

type gameData struct {
	GameID   string `json:"game-id"`
	Username string `json:"username"`
}

type errData struct {
	Err  string `json:"error"`
	Desc string `json:"description,omitempty"`
}

//The following two functions could _definitely_ be in main at their base functionality, but they are here in case we want to do any additional validation in the future.
func upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func clientResponse(response *response, socket *websocket.Conn) (error, error) {
	r, marshalErr := json.Marshal(response)
	writeErr := socket.WriteMessage(1, r)
	return marshalErr, writeErr
}
