package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"spyfall/db"
	"time"

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

type clientmessage struct {
	Trigger, Data string
}

type gamejoin struct {
	Code, Username string
}

//Configurations

var dbaddr = "127.0.0.1"
var dbport = "27017"
var dbname = "spyfall"
var dbgamecollection = "games"

var apiport = "8080"

func main() {
	print("general", "Attempting to connect to database: \""+dbname+"\" at: "+dbaddr+":"+dbport)

	dbo := db.DBO{
		Server:         dbaddr + ":" + dbport,
		Database:       dbname,
		GameCollection: dbgamecollection,
	}

	err := db.Connect(dbo)
	if err != nil {
		fmt.Println(err)
		print("db", "Unable to connect to database!")
	} else {
		print("db", "Connected to database!")
	}

	print("general", "Starting web server on port: "+apiport)
	http.HandleFunc("/", handler)
	http.HandleFunc("/api", api)
	http.ListenAndServe(":"+apiport, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("../../../public"))
}

func api(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	print("ws", "Client Connected!")

	for {
		var clientMessage clientmessage
		connection.ReadJSON(&clientMessage)

		print("ws", "Trigger: "+clientMessage.Trigger)

		switch clientMessage.Trigger {
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
		default:
			clientReturn("ERROR", "Not a valid trigger")
		}

	}
}

func createGame(data string) []byte {
	var gameJoin gamejoin
	var code string
	json.Unmarshal([]byte(data), &gameJoin)

	if gameJoin.Code == "" {
		code = generateCode()

		err := db.NewGame(code, "spy-school")
		if err != nil {
			fmt.Println(err)
		}

		print("api", "Game \""+code+"\" doesn't exist, creating...")

		returnMessage := clientReturn("OK", code)
		return returnMessage
	} else {
		code = gameJoin.Code

		err := db.NewGame(code, "spy-school")
		if err != nil {
			fmt.Println(err)
		}

		if false { //Check if game is already in db
			print("api", "Game \""+code+"\" already exists in the database, error")
			return clientReturn("ERROR", "Game code: "+code+" already exists")
		} else {
			print("api", "Game \""+code+"\" doesn't exist, creating...")
			return clientReturn("OK", code)
		}
	}
}

func joinGame(data string) []byte {
	var gameJoin gamejoin
	json.Unmarshal([]byte(data), &gameJoin)

	if gameJoin.Code == "" {
		print("api", "Game code blank, error")
		return clientReturn("ERROR", "Game code cannot be blank")
	} else {
		err := db.AddPlayer(gameJoin.Code, gameJoin.Username)

		if err != nil { //Check if game is in db
			print("api", "Game \""+gameJoin.Code+"\" not found in database, error")
			return clientReturn("ERROR", "Game code: \""+gameJoin.Code+"\" doesn't exist")
		}
		print("api", "Game \""+gameJoin.Code+"\" found in database, joining...")
		return clientReturn("OK", gameJoin.Code)
	}
}

func clientReturn(returnResponse, returnData string) []byte {
	rd, err := json.Marshal(&clientmessage{
		Trigger: returnResponse,
		Data:    returnData,
	})
	if err != nil {
		fmt.Println(err)
	}
	print("ws", "Returning to client: "+string(rd))
	return rd
}

func generateCode() string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyz")
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, 6)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	print("general", "generateCode() called, code: "+string(b))
	return string(b)
}

func print(loglevel, message string) {
	var prefix string
	switch loglevel {
	case "api":
		prefix = "[API] "
	case "general":
		prefix = "[General] "
	case "db":
		prefix = "[Database] "
	case "ws":
		prefix = "[Websockets] "
	}

	fmt.Println(prefix + message)
}