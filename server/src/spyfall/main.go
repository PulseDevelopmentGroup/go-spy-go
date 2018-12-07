package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"spyfall/db"
	"spyfall/websockets"
	"time"
)

var configFile = "config.json"

type Config struct {
	Databse struct {
		Host       string `json:"host"`
		Port       string `json:"port"`
		Name       string `json:"name"`
		Collection string `json:"collection"`
	} `json:"database"`
	Api struct {
		Port string `json:"port"`
	} `json:"api"`
	Locations []string `json:"locations"`
}

type Location struct {
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
	Spies int      `json:"spies"`
}

type Locations struct {
	Locations []Location `json:"locations"`
}

func main() {
	config, err := readConfig(configFile)
	if err != nil {
		fmt.Println(err)
	}

	print("db", "Attempting to connect to database: \""+config.Databse.Name+"\" at: "+config.Databse.Host+":"+config.Databse.Port)

	db.Connect(&db.DBO{
		Server:         config.Databse.Host + ":" + config.Databse.Port,
		Database:       config.Databse.Name,
		GameCollection: config.Databse.Collection,
	})
	if err != nil {
		fmt.Println(err)
		print("db", "Unable to connect to database!")
	} else {
		print("db", "Connected to database!")
	}

	print("general", "Starting web server on port: "+config.Api.Port)
	http.HandleFunc("/", handler)
	http.HandleFunc("/api", apiHandler)
	http.ListenAndServe(":"+config.Api.Port, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("public"))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	connection, err := websockets.Upgrade(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		var request websockets.Request
		var data websockets.GameData
		connection.ReadJSON(&request)

		json.Unmarshal([]byte(request.Data), &data)

		print("ws", "Recieved: Kind: "+request.Kind+" Data: "+request.Data)

		switch request.Kind {
		case "CREATE_GAME":
			createGame(data, connection)
		case "JOIN_GAME":
			joinGame(data, connection)
		case "START_GAME":
			//websockets.ClientResponse(startGame(request.Data), connection)
			//Start game
		case "STOP_GAME":
			//Stop game
		case "LEAVE_GAME":
			//Leave game
		default:
			websockets.Send(&websockets.Response{
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

func createGame(gameData websockets.GameData, connection websockets.Connection) {
	var code string

	if gameData.GameId == "" {
		code = generateCode()
		gameData.GameId = code
	} else {
		code = gameData.GameId
	}

	err := db.NewGame(code, "spy-school") //This needs to be randomly generated

	if err != nil {
		fmt.Println(err)
		if err == fmt.Errorf("GAME_ALREADY_EXISTS") {
			print("api", "Game \""+code+"\" already exists in database.")
			websockets.Send(&websockets.Response{
				Kind: "CREATE_GAME",
				Data: marshal(&websockets.GameData{
					GameId:   code,
					Username: gameData.Username,
				}),
				Err: marshal(&websockets.ErrData{
					Err:  "GAME_ALREADY_EXISTS",
					Desc: "Game: \"" + code + "\" already exists in database.",
				}),
			}, connection)
			return
		} else {
			print("api", "There was a big problem")
			websockets.Send(&websockets.Response{
				Kind: "CREATE_GAME",
				Data: marshal(&websockets.GameData{
					GameId:   code,
					Username: gameData.Username,
				}),
				Err: marshal(&websockets.ErrData{
					Err:  "UNKNOWN_ERROR",
					Desc: "See the server log",
				}),
			}, connection)
			return
		}
	}
	print("api", "Game \""+code+"\" doesn't exist, creating...")
	websockets.Send(&websockets.Response{
		Kind: "CREATE_GAME",
		Data: marshal(&websockets.GameData{
			GameId:   code,
			Username: gameData.Username,
		}),
	}, connection)
	joinGame(gameData, connection)
	return
}

func joinGame(gameData websockets.GameData, connection websockets.Connection) {
	if gameData.GameId == "" {
		print("api", "Game code blank, error")
		websockets.Send(&websockets.Response{
			Kind: "JOIN_GAME",
			Data: marshal(gameData),
			Err: marshal(&websockets.ErrData{
				Err:  "NO_GAME_CODE",
				Desc: "Good luck joining a game with no code!",
			}),
		}, connection)
		return
	}

	pid, err := db.AddPlayer(gameData.GameId, gameData.Username)
	if err != nil {
		fmt.Println(err)
		print("api", "Error: "+err.Error())

		if err.Error() == "NO_GAME_EXISTS" {
			websockets.Send(&websockets.Response{
				Kind: "JOIN_GAME",
				Data: marshal(gameData),
				Err: marshal(&websockets.ErrData{
					Err:  "NO_GAME",
					Desc: "The game: \"" + gameData.GameId + "\" does not exist.",
				}),
			}, connection)
			return
		}
		if err.Error() == "USER_ALREADY_EXISTS" {
			websockets.Send(&websockets.Response{
				Kind: "JOIN_GAME",
				Data: marshal(gameData),
				Err: marshal(&websockets.ErrData{
					Err:  "DUP_USER",
					Desc: "A user with the username: \"" + gameData.Username + "\" already exists in game: \"" + gameData.GameId + "\".",
				}),
			}, connection)
			return
		}
		websockets.Send(&websockets.Response{
			Kind: "JOIN_GAME",
			Data: marshal(gameData),
			Err: marshal(&websockets.ErrData{
				Err:  "ERROR",
				Desc: "This shouldn't happen, see the server log for details.",
			}),
		}, connection)
		return
	}
	print("api", "Game \""+gameData.GameId+"\" found in database, joining...")
	websockets.Clients[pid] = connection
	websockets.Send(&websockets.Response{
		Kind: "JOIN_GAME",
		Data: marshal(gameData),
	}, connection)
	return
}

/*func getLocations(file string) (Locations, error) {
	var locations Locations
	osFile, err := os.Open(file)
	defer osFile.Close()
	if err != nil {
		return locations, err
	}

	json.NewDecoder(osFile).Decode(&locations)
	return locations, err
}*/

func readConfig(file string) (Config, error) {
	var config Config
	osFile, err := os.Open(file)
	defer osFile.Close()
	if err != nil {
		return config, err
	}

	json.NewDecoder(osFile).Decode(&config)
	fmt.Println(config)
	return config, err
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

func marshal(input interface{}) string {
	r, err := json.Marshal(input)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(r)
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
