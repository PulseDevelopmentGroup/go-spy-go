package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"spyfall/db"
	"spyfall/websockets"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
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

var config Config

func main() {
	rConfig, err := readConfig(configFile)
	if err != nil {
		fmt.Println(err)
	}
	config = rConfig

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
		var gameData websockets.GameData

		err := connection.ReadJSON(&request)
		if err != nil {
			connection.Close()
			print("ws", "Websockets connection terminated")
			return
		}

		json.Unmarshal([]byte(request.Data), &gameData)

		print("ws", "Recieved: Kind: "+request.Kind+" Data: "+request.Data)

		switch request.Kind {
		case "CREATE_GAME":
			createGame(gameData, connection)
		case "JOIN_GAME":
			joinGame(gameData, connection)
		case "START_GAME":
			startGame(connection)
		case "STOP_GAME":
			//Stop game
		case "LEAVE_GAME":
			leaveGame(connection)
		default:
			websockets.SendToPlayer(&websockets.Response{
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

//Main game Functions
func createGame(gameData websockets.GameData, connection *websocket.Conn) {
	var code string

	if gameData.GameID == "" {
		code = generateCode()
		gameData.GameID = code
	} else {
		code = gameData.GameID
	}

	err := db.AddGame(code)

	if err != nil {
		fmt.Println(err)
		if err == fmt.Errorf("GAME_ALREADY_EXISTS") {
			print("api", "Game \""+code+"\" already exists in database.")
			websockets.SendToPlayer(&websockets.Response{
				Kind: "CREATE_GAME",
				Data: marshal(&websockets.GameData{
					GameID:   code,
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
			websockets.SendToPlayer(&websockets.Response{
				Kind: "CREATE_GAME",
				Data: marshal(&websockets.GameData{
					GameID:   code,
					Username: gameData.Username,
				}),
				Err: marshal(&websockets.ErrData{
					Err:  "UNKNOWN_ERROR",
					Desc: "This shouldn't happen, see the server log for details.",
				}),
			}, connection)
			return
		}
	}
	print("api", "Game \""+code+"\" doesn't exist, creating...")
	websockets.SendToPlayer(&websockets.Response{
		Kind: "CREATE_GAME",
		Data: marshal(&websockets.GameData{
			GameID:   code,
			Username: gameData.Username,
		}),
	}, connection)
	joinGame(gameData, connection)
	return
}

func joinGame(gameData websockets.GameData, connection *websocket.Conn) {
	if gameData.GameID == "" {
		print("api", "Game code blank, error")
		websockets.SendToPlayer(&websockets.Response{
			Kind: "JOIN_GAME",
			Data: marshal(gameData),
			Err: marshal(&websockets.ErrData{
				Err:  "NO_GAME_CODE",
				Desc: "Good luck joining a game with no code!",
			}),
		}, connection)
		return
	}

	pid, err := db.AddPlayer(gameData.GameID, gameData.Username)
	if err != nil {
		fmt.Println(err)
		print("api", "Error: "+err.Error())

		if err.Error() == "NO_GAME_EXISTS" {
			websockets.SendToPlayer(&websockets.Response{
				Kind: "JOIN_GAME",
				Data: marshal(gameData),
				Err: marshal(&websockets.ErrData{
					Err:  "NO_GAME",
					Desc: "The game: \"" + gameData.GameID + "\" does not exist.",
				}),
			}, connection)
			return
		}
		if err.Error() == "USER_ALREADY_EXISTS" {
			websockets.SendToPlayer(&websockets.Response{
				Kind: "JOIN_GAME",
				Data: marshal(gameData),
				Err: marshal(&websockets.ErrData{
					Err:  "USER_ALREADY_EXISTS",
					Desc: "A user with the username: \"" + gameData.Username + "\" already exists in game: \"" + gameData.GameID + "\".",
				}),
			}, connection)
			return
		}
		print("api", "There was a big problem")
		websockets.SendToPlayer(&websockets.Response{
			Kind: "JOIN_GAME",
			Data: marshal(gameData),
			Err: marshal(&websockets.ErrData{
				Err:  "UNKNOWN_ERROR",
				Desc: "This shouldn't happen, see the server log for details.",
			}),
		}, connection)
		return
	}
	print("api", "Game \""+gameData.GameID+"\" found in database, joining...")

	websockets.ClientById[pid.Hex()] = connection
	websockets.ClientByConn[connection] = pid.Hex()

	print("ws", "Associating this connection with player id: "+pid.Hex())

	websockets.SendToPlayer(&websockets.Response{
		Kind: "JOIN_GAME",
		Data: marshal(gameData),
	}, connection)

	return
}

func startGame(connection *websocket.Conn) {
	gamecode, err := db.GetGameCode(websockets.ClientByConn[connection])
	if err != nil {
		fmt.Println(err)
		return
	}

	print("game", "Starting game: "+gamecode)

	locations, err := getLocations(config.Locations[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	rand.Seed(time.Now().UnixNano())
	location := locations.Locations[rand.Intn(len(locations.Locations))]
	print("game", "--------------------------------------------------------------------")
	print("game", "Location: "+location.Name)
	print("game", "Roles: ["+strings.Join(location.Roles, ", ")+"]")
	print("game", "Number of Spies: "+strconv.Itoa(location.Spies))
	print("game", "--------------------------------------------------------------------")

	db.SetLocation(gamecode, location.Name)

	setRoles(gamecode, location.Roles, location.Spies)
}

func leaveGame(connection *websocket.Conn) {
	pid := websockets.ClientByConn[connection]
	player, err := db.GetPlayer(pid)
	if err != nil {
		fmt.Println(err)
		return
	}
	gamecode, err := db.GetGameCode(pid)
	if err != nil {
		fmt.Println(err)
		return
	}
	game, err := db.GetGameData(gamecode)
	if err != nil {
		fmt.Println(err)
		return
	}
	print("game", "Removing player: "+player.Username+" ("+pid+") from game: "+game.GameCode)

	for i := 0; i < len(game.Players); i++ {
		websockets.SendToPlayer(&websockets.Response{
			Kind: "LEAVE_GAME",
			Data: marshal(&websockets.LeaveData{
				Username: player.Username,
				Reason:   "Requested to leave/lost connection",
			}),
		}, websockets.ClientById[game.Players[i].PlayerID.Hex()])
	}

	db.DelPlayer(pid)
	delete(websockets.ClientByConn, connection)
	delete(websockets.ClientById, pid)
}

//Helper Game functions
func setRoles(gamecode string, inputRoles []string, spies int) {
	players, err := db.GetPlayers(gamecode)
	if err != nil {
		fmt.Println(err)
		return
	}
	print("game", "There are "+strconv.Itoa(len(players))+" players who need roles assigned. This location calls for "+strconv.Itoa(spies)+" spy/spies")

	//Clear roles and spies
	//This should be moved to stop game function eventually
	for i := 0; i < len(players); i++ {
		print("game", "Clearing roles for: "+players[i].Username+" ("+players[i].PlayerID.Hex()+")")
		players[i].Role = "null"
		players[i].Spy = false
	}

	for i := 0; i < spies; i++ {
		rand.Seed(time.Now().UnixNano())
		spyIndex := rand.Intn(len(players))
		if players[spyIndex].Spy {
			print("game", players[i].Username+" ("+players[i].PlayerID.Hex()+") is already a spy! Skipping.")
		} else {
			print("game", "A spy is: "+players[spyIndex].Username+" ("+players[spyIndex].PlayerID.Hex()+")")
			players[spyIndex].Role = "spy"
			players[spyIndex].Spy = true

			db.SetPlayer(&players[spyIndex])
		}
	}

	roles := inputRoles
	for i := 0; i < len(players); i++ {
		rand.Seed(time.Now().UnixNano())
		if players[i].Spy {
			print("game", "Player: "+players[i].Username+" ("+players[i].PlayerID.Hex()+") is a spy! They cannot have a role!")
		} else {
			if len(roles) <= 1 {
				roles = inputRoles
				print("game", "All roles assigned! Resetting role list!")
			}
			roleIndex := rand.Intn(len(roles))
			players[i].Role = roles[roleIndex]
			print("game", "Assigning role: "+players[i].Role+" to player: "+players[i].Username+" ("+players[i].PlayerID.Hex()+")")
			db.SetPlayer(&players[i])

			roles = append(roles[:roleIndex], roles[roleIndex+1:]...)
		}
	}
}

func getLocations(file string) (Locations, error) {
	var locations Locations

	osFile, err := os.Open(file)
	defer osFile.Close()
	if err != nil {
		return locations, err
	}

	json.NewDecoder(osFile).Decode(&locations)
	return locations, err
}

func generateCode() string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyz")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 6)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	print("game", "generateCode() called, code: "+string(b))
	return string(b)
}

//Utility functions
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
		prefix = "] [   API    ] "
	case "general":
		prefix = "] [   Main   ] "
	case "db":
		prefix = "] [ Database ] "
	case "ws":
		prefix = "] [Websockets] "
	case "game":
		prefix = "] [   Game   ] "
	}
	fmt.Println("[" + time.Now().Format("01-02-2006 15:04:05") + prefix + message)
}
