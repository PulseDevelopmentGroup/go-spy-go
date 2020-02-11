package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

//Configurations

var dbaddr = "127.0.0.1"
var dbport = "27017"
var dbname = "spyfall"
var dbgamecollection = "games"

var apiport = "8080"

func main() {
	print("general", "Attempting to connect to database: \""+dbname+"\" at: "+dbaddr+":"+dbport)

	dbo := &dbo{
		Server:         dbaddr + ":" + dbport,
		Database:       dbname,
		GameCollection: dbgamecollection,
	}

	err := connect(dbo)
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
	connection, err := upgrade(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		var request request
		connection.ReadJSON(&request)

		fmt.Println(request)

		print("ws", "Recieved: Kind: "+request.Kind+" Data: "+request.Data)

		switch request.Kind {
		case "CREATE_GAME":
			clientResponse(createGame(request.Data), connection)
		case "JOIN_GAME":
			clientResponse(joinGame(request.Data), connection)
		case "START_GAME":
			//Start game
		case "STOP_GAME":
			//Stop game
		case "LEAVE_GAME":
			//Leave game
		default:
			fmt.Errorf("\"%s\" is not a valid kind", request.Kind)
			clientResponse(&response{
				Kind: request.Kind,
				Data: request.Data,
				Err: marshal(&errData{
					Err:  "NOT_VALID_KIND",
					Desc: "\"" + request.Kind + "\" is not a valid kind.",
				}),
			}, connection)
		}
	}
}

func createGame(data string) *response {
	var code string
	var gameInfo gameData
	json.Unmarshal([]byte(data), &gameInfo)

	if gameInfo.GameID == "" {
		code = generateCode()
	} else {
		code = gameInfo.GameID
	}

	err := newGame(code, "spy-school") //This needs to be randomly generated
	if err != nil {
		fmt.Println(err)
		print("api", "Game \""+code+"\" already exists in database.")
		return &response{
			Kind: "CREATE_GAME",
			Data: marshal(&gameData{
				GameID:   code,
				Username: gameInfo.Username,
			}),
			Err: marshal(&errData{
				Err:  "GAME_ALREADY_EXISTS",
				Desc: "Game: \"" + code + "\" already exists in database.",
			}),
		}
	}

	print("api", "Game \""+code+"\" doesn't exist, creating...")
	joinGame(marshal(&gameData{
		GameID:   code,
		Username: gameInfo.Username,
	}))
	return &response{
		Kind: "CREATE_GAME",
		Data: marshal(&gameData{
			GameID:   code,
			Username: gameInfo.Username,
		}),
	}
}

func joinGame(data string) *response {
	var gameData gameData
	json.Unmarshal([]byte(data), &gameData)

	if gameData.GameID == "" {
		print("api", "Game code blank, error")
		return &response{
			Kind: "JOIN_GAME",
			Data: data,
			Err: marshal(&errData{
				Err:  "NO_GAME_CODE",
				Desc: "Good luck joining a game with no code!",
			}),
		}
	}

	err := addPlayer(gameData.GameID, gameData.Username)
	if err != nil {
		fmt.Println(err)
		print("api", "Game \""+gameData.GameID+"\" not found in database, error")
		return &response{
			Kind: "JOIN_GAME",
			Data: data,
			Err: marshal(&errData{
				Err:  "NO_GAME",
				Desc: "The game: \"" + gameData.GameID + "\" does not exist.",
			}),
		}

	}
	print("api", "Game \""+gameData.GameID+"\" found in database, joining...")
	return &response{
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

func marshal(input interface{}) string {
	r, err := json.Marshal(input)
	if err != nil {
		fmt.Println(err.Error)
	}
	return string(r)
}
