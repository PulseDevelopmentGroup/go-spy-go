package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"spyfall/db"
	"spyfall/websockets"
	"time"
)

//Configurations

var dbaddr = "127.0.0.1"
var dbport = "27017"
var dbname = "spyfall"
var dbgamecollection = "games"

var apiport = "8080"

func main() {
	print("db", "Attempting to connect to database: \""+dbname+"\" at: "+dbaddr+":"+dbport)

	dbo := &db.DBO{
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
	http.FileServer(http.Dir("../../../public")) //Should probably fix this sometime
}

func api(w http.ResponseWriter, r *http.Request) {
	connection, err := websockets.Upgrade(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		var request websockets.Request
		connection.ReadJSON(&request)

		fmt.Println(request)

		print("ws", "Recieved: Kind: "+request.Kind+" Data: "+request.Data)

		switch request.Kind {
		case "CREATE_GAME":
			websockets.ClientResponse(createGame(request.Data), connection)
		case "JOIN_GAME":
			websockets.ClientResponse(joinGame(request.Data), connection)
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

func createGame(data string) *websockets.Response {
	var code string
	var gameData websockets.GameData
	json.Unmarshal([]byte(data), &gameData)

	if gameData.GameId == "" {
		code = generateCode()
	} else {
		code = gameData.GameId
	}

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
	} else {
		print("api", "Game \""+code+"\" doesn't exist, creating...")
		joinGame(marshal(&websockets.GameData{
			GameId:   code,
			Username: gameData.Username,
		}))
		return &websockets.Response{
			Kind: "CREATE_GAME",
			Data: marshal(&websockets.GameData{
				GameId:   code,
				Username: gameData.Username,
			}),
		}
	}
}

func joinGame(data string) *websockets.Response {
	var gameData websockets.GameData
	json.Unmarshal([]byte(data), &gameData)

	if gameData.GameId == "" {
		print("api", "Game code blank, error")
		return &websockets.Response{
			Kind: "JOIN_GAME",
			Data: data,
			Err: marshal(&websockets.ErrData{
				Err:  "NO_GAME_CODE",
				Desc: "Good luck joining a game with no code!",
			}),
		}
	}

	err := db.AddPlayer(gameData.GameId, gameData.Username)
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
	fmt.Println(prefix + message)
}

func marshal(input interface{}) string {
	r, err := json.Marshal(input)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(r)
}
