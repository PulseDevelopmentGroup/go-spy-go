package db

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DBO struct {
	Server         string
	Database       string
	GameCollection string
}

type Player struct {
	PlayerID bson.ObjectId `bson:"playerid" json:"playerid"`
	Username string        `bson:"username" json:"username"`
	Role     string        `bson:"role" json:"role"`
	Spy      bool          `bson:"spy" json:"spy"`
}

type Game struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	GameCode string        `bson:"gamecode" json:"gamecode"`
	Location string        `bson:"location" json:"location"`
	Players  []Player      `bson:"players" json:"players"`
}

var collection *mgo.Collection

//Connect connects to the database with the supplied database object
func Connect(dbo *DBO) error {
	session, err := mgo.Dial(dbo.Server)
	if err != nil {
		return err
	}

	db := session.DB(dbo.Database)
	collection = db.C(dbo.GameCollection)

	return nil
}

//AddPlayer adds a new player to the database
//Returns the player id and an error
func AddPlayer(gamecode, username string) (bson.ObjectId, error) {
	err := checkGame(gamecode)
	if err != nil {
		if err == fmt.Errorf("NO_GAME_EXISTS") {
			return "", fmt.Errorf("NO_GAME_EXISTS")
		}
		return "", err
	}
	usernameCount, err := collection.Find(bson.M{"gamecode": gamecode, "players": bson.M{"$elemMatch": bson.M{"username": username}}}).Limit(1).Count()
	if err != nil {
		return "", err
	}
	if usernameCount > 0 {
		return "", fmt.Errorf("USER_ALREADY_EXISTS")
	}
	pid := bson.NewObjectId()
	udErr := collection.Update(bson.M{"gamecode": gamecode}, bson.M{"$push": bson.M{"players": &Player{
		PlayerID: pid,
		Username: username,
		Role:     "null",
		Spy:      false,
	}}})
	return pid, udErr
}

//GetPid returns the id of the player in a specific game based on their username
func GetPid(gamecode, username string) (bson.ObjectId, error) {
	players, err := GetPlayers(gamecode)
	if err != nil {
		fmt.Println(err)
	}
	for i := range players {
		if players[i].Username == username {
			fmt.Println("GetPid - " + players[i].Username)
			return players[i].PlayerID, err
		}
	}
	return bson.NewObjectId(), fmt.Errorf("NO_PLAYER_EXISTS")
}

//GetPlayers returns an array of players in a game
func GetPlayers(gamecode string) ([]Player, error) {
	game, gdErr := GetGameData(gamecode)
	if gdErr != nil {
		return nil, gdErr
	}
	return game.Players, nil
}

//GetPlayer returns a player struct with the data of a player matching a username
//TODO: TEST THIS
//TODO: Add error handling to this
func GetPlayer(pid string) Player {
	game := Game{}
	collection.Find(bson.M{"players.playerid": bson.ObjectIdHex(pid)}).Select(bson.M{"players.$": 1}).One(&game)
	return game.Players[0]
}

//SetPlayer updates a player's record with elements from a new supplied player
//When supplying the "Player" variable, don't change it's object id. You're gonna have a bad time
//TODO: TEST THIS
func SetPlayer(player *Player) error {
	updatePlayer := GetPlayer(player.PlayerID.Hex())

	if updatePlayer.Username != player.Username {
		updatePlayer.Username = player.Username
	}
	if updatePlayer.Role != player.Role {
		updatePlayer.Role = player.Role
	}
	if updatePlayer.Spy != player.Spy {
		updatePlayer.Spy = player.Spy
	}

	collection.Update(bson.M{"players.playerid": player.PlayerID}, bson.M{"$set": bson.M{"players.$": updatePlayer}})
	return nil
}

//SetLocation sets location of a specific game
func SetLocation(gamecode, location string) error {
	err := checkGame(gamecode)
	if err != nil {
		if err == fmt.Errorf("NO_GAME_EXISTS") {
			return fmt.Errorf("NO_GAME_EXISTS")
		}
		return err
	}
	collection.Update(bson.M{"gamecode": gamecode}, bson.M{"$set": bson.M{"location": location}})

	return nil
}

//GetLocation returns a string of the location of a specific game
func GetLocation(gamecode string) (string, error) {
	game, gdErr := GetGameData(gamecode)
	if gdErr != nil {
		if gdErr == fmt.Errorf("NO_GAME_EXISTS") {
			return "", fmt.Errorf("NO_GAME_EXISTS")
		}
		return "", gdErr
	}

	return game.Location, nil
}

//NewGame creates a new game in the DB with the provided gamecode
//Returns error if game already exists or cannot be created
//NOTE: LOCATION NOT SET HERE
func NewGame(gamecode string) error {
	return insertGame(&Game{
		ID:       bson.NewObjectId(),
		GameCode: gamecode,
		Location: "null",
		Players:  nil,
	})
}

//GetGameData returns the Game struct that matches the provided gamecode
//Returns error if game does not exist
//TODO Test this more
func GetGameData(gamecode string) (Game, error) {
	checkGameErr := checkGame(gamecode)
	result := Game{}
	if checkGameErr != nil {
		return result, checkGameErr
	}
	err := collection.Find(bson.M{"gamecode": gamecode}).One(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

//Returns an error if game doesn't exist. Returns nothing if a game exists
func checkGame(gamecode string) error {
	gameCount, err := collection.Find(bson.M{"gamecode": gamecode}).Limit(1).Count()
	if err != nil {
		return err
	}
	if gameCount == 0 {
		return fmt.Errorf("NO_GAME_EXISTS")
	}
	return nil
}

//Adds a game to the DB
func insertGame(entry *Game) error {
	count, err := collection.Find(bson.M{"gamecode": entry.GameCode}).Limit(1).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("GAME_ALREADY_EXISTS")
	}
	return collection.Insert(entry)
}

//SetSpies sets the spies in a game with an array of usernames
//Depreciated in favor of the SetPlayer function
func SetSpies(gamecode string, spies []string) error {
	err := checkGame(gamecode)
	if err != nil {
		if err == fmt.Errorf("NO_GAME_EXISTS") {
			return fmt.Errorf("NO_GAME_EXISTS")
		}
		return err
	}
	for i := 0; i < len(spies); i++ {
		collection.Update(bson.M{"gamecode": gamecode, "players.username": spies[i]}, bson.M{"$set": bson.M{"players.$.spy": true}})
	}

	return nil
}

//GetSpies returns an array of strings with the usernames of the game's spies
//Depreciated in favor of the GetPlayers function
func GetSpies(gamecode string) ([]string, error) {
	game, gdErr := GetGameData(gamecode)
	spies := []string{}
	if gdErr != nil {
		if gdErr == fmt.Errorf("NO_GAME_EXISTS") {
			return nil, fmt.Errorf("NO_GAME_EXISTS")
		}
		return nil, gdErr
	}
	for i := 0; i < len(game.Players); i++ {
		if game.Players[i].Spy {
			spies = append(spies, game.Players[i].Username)
		}
	}

	return spies, nil
}
