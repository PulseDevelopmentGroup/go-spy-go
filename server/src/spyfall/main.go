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
	Kind, Data string
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
	print("db", "Attempting to connect to database: \""+dbname+"\" at: "+dbaddr+":"+dbport)

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

		print("ws", "Kind: "+clientMessage.Kind)

		switch clientMessage.Kind {
		case "CREATE_GAME":
			connection.WriteMessage(1, createGame(clientMessage.Data))
		case "JOIN_GAME":
			connection.WriteMessage(1, joinGame(clientMessage.Data))
		case "START_GAME":
			fmt.Println(db.GetLocation("wyquut"))
		case "STOP_GAME":
			fmt.Println(db.SetLocation("wyquut", "new-location"))
		case "LEAVE_GAME":
			//Leave game
		default:
			websockets.ClientResponse(&websockets.Response{
				Kind: request.Kind,
				Data: request.Data,
				Err: marshal(&websockets.ErrData{
					Err:  "NOT_VALID_KIND",
					Desc: "\"" + request.Kind + "\" is not a valid kind.",
				}),
			}, connection)
		}

	}
}

func createGame(data string) []byte {
	var gameJoin gamejoin
	var code string
	json.Unmarshal([]byte(data), &gameJoin)

	if gameJoin.Code == "" {
		code = generateCode()

	err := db.NewGame(code, "spy-school") //This needs to be randomly generated

	if err != nil {
		fmt.Println(err)
		if err == fmt.Errorf("GAME_ALREADY_EXISTS") {
			print("api", "Game \""+code+"\" already exists in database.")
			return &websockets.Response{
				Kind: "CREATE_GAME",
				Data: marshal(&websockets.GameData{
					GameId:   code,
					Username: gameData.Username,
				}),
				Err: marshal(&websockets.ErrData{
					Err:  "GAME_ALREADY_EXISTS",
					Desc: "Game: \"" + code + "\" already exists in database.",
				}),
			}
		} else {
			print("api", "There was a big problem")
			return &websockets.Response{
				Kind: "CREATE_GAME",
				Data: marshal(&websockets.GameData{
					GameId:   code,
					Username: gameData.Username,
				}),
				Err: marshal(&websockets.ErrData{
					Err:  "UNKNOWN_ERROR",
					Desc: "See the server log",
				}),
			}
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
		Kind: returnResponse,
		Data: returnData,
	})
	if err != nil {
		fmt.Println(err)
		print("api", "Error: "+err.Error())

		if err.Error() == "NO_GAME_EXISTS" {
			return &websockets.Response{
				Kind: "JOIN_GAME",
				Data: data,
				Err: marshal(&websockets.ErrData{
					Err:  "NO_GAME",
					Desc: "The game: \"" + gameData.GameId + "\" does not exist.",
				}),
			}
		}
		if err.Error() == "USER_ALREADY_EXISTS" {
			return &websockets.Response{
				Kind: "JOIN_GAME",
				Data: data,
				Err: marshal(&websockets.ErrData{
					Err:  "DUP_USER",
					Desc: "A user with the username: \"" + gameData.Username + "\" already exists in game: \"" + gameData.GameId + "\".",
				}),
			}
		}

		return &websockets.Response{
			Kind: "JOIN_GAME",
			Data: data,
			Err: marshal(&websockets.ErrData{
				Err:  "ERROR",
				Desc: "This shouldn't happen, see the server log for details.",
			}),
		}

	}
	print("api", "Game \""+gameData.GameId+"\" found in database, joining...")

	return &websockets.Response{
		Kind: "JOIN_GAME",
		Data: data,
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
		prefix = "[   API    ] "
	case "general":
		prefix = "[ General  ] "
	case "db":
		prefix = "[ Database ] "
	case "ws":
		prefix = "[Websockets] "
	}

func marshal(input interface{}) string {
	r, err := json.Marshal(input)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(r)
}
