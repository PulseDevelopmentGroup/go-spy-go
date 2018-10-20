package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ClientMessage struct {
	Type, Data string
}

type GameJoin struct {
	Code, Username string
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/api", api)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("public"))
}

func api(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Client Connected!")

	for {
		var clientMessage ClientMessage
		connection.ReadJSON(&clientMessage)

		switch clientMessage.Type {
		case "create-game":
			connection.WriteMessage(1, createGame(clientMessage.Data))
		case "join-game":
			connection.WriteMessage(1, joinGame(clientMessage.Data))
		case "start-game":
			//Start game
		case "stop-game":
			//Stop game
		case "leave-game":
			//Leave game
		}

	}
}

func createGame(data string) []byte {
	var gameJoin GameJoin
	var code string
	json.Unmarshal([]byte(data), &gameJoin)

	fmt.Println(data)

	if gameJoin.Code == "" {
		code = generateCode()
		//Add game to database (if it doesn't already exist)
		//Inform user of generated game w/code

		return clientReturn("OK", code)
	} else {
		code = gameJoin.Code
		//Add game to database (if it doesn't already exist)
		//Inform user of generated game w/code

		if false { //Check if game is already in db
			return clientReturn("ERROR", "Game code: "+code+" already exists")
		} else {
			return clientReturn("OK", code)
		}
	}
}

func joinGame(data string) []byte {
	var gameJoin GameJoin
	json.Unmarshal([]byte(data), &gameJoin)

	if gameJoin.Code == "" {
		return clientReturn("ERROR", "Game code cannot be blank")
	} else {
		if false { //Check if game is in db
			return clientReturn("OK", "Joining game code: "+gameJoin.Code)
		} else {
			return clientReturn("ERROR", "Game code: "+gameJoin.Code+" doesn't exist")
		}
	}
}

func clientReturn(returnType, data string) []byte {
	returnData, err := json.Marshal(&ClientMessage{
		Type: returnType,
		Data: data,
	})
	if err != nil {
		fmt.Println(err)
	}
	return returnData
}

func generateCode() string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyz")
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, 6)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
