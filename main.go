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
	response, data string
}

type GameJoin struct {
	code, username string
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

		switch clientMessage.response {
		case "create-game":
			fmt.Println("Create Game")
			connection.WriteMessage(1, createGame(clientMessage.data))
		case "join-game":
			fmt.Println("Join Game")
			connection.WriteMessage(1, joinGame(clientMessage.data))
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

	if gameJoin.code == "" {
		code = generateCode()
		//Add game to database (if it doesn't already exist)
		//Inform user of generated game w/code

		returnMessage := clientReturn("OK", code)
		fmt.Println(string(returnMessage))

		return returnMessage
	} else {
		code = gameJoin.code
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

	if gameJoin.code == "" {
		return clientReturn("ERROR", "Game code cannot be blank")
	} else {
		if false { //Check if game is in db
			return clientReturn("OK", gameJoin.code)
		} else {
			return clientReturn("ERROR", "Game code: \""+gameJoin.code+"\" doesn't exist")
		}
	}
}

func clientReturn(returnResponse, returnData string) []byte {
	rd, err := json.Marshal(&ClientMessage{
		response: returnResponse,
		data:     returnData,
	})
	if err != nil {
		fmt.Println(err)
	}
	return rd
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
