package db

import (
	"fmt"
	"reflect"

	mgo "gopkg.in/mgo.v2"
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
	Active   bool          `bson:"active" json:"active"`
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
	game, err := GetGame(gamecode)
	if err != nil {
		return "", err
	}

	if game.Active {
		return "", fmt.Errorf("GAME_IN_PROGRESS")
	}
	for i := 0; i < len(game.Players); i++ {
		if game.Players[i].Username == username {
			return "", fmt.Errorf("USER_ALREADY_EXISTS")
		}
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
			return players[i].PlayerID, err
		}
	}
	return bson.NewObjectId(), fmt.Errorf("NO_PLAYER_EXISTS")
}

//GetPlayer returns a player struct with the data of a player matching a pid
func GetPlayer(pid string) (Player, error) {
	game := Game{}
	err := collection.Find(bson.M{"players.playerid": bson.ObjectIdHex(pid)}).Select(bson.M{"players.$": 1}).One(&game)
	return game.Players[0], err
}

//SetPlayer updates a player's record with elements from a new supplied player
//When supplying the "Player" variable, don't change it's object id. You're gonna have a bad time
//TODO: TEST THIS
func SetPlayer(player *Player) error {
	updatePlayer, err := GetPlayer(player.PlayerID.Hex())
	if err != nil {
		return err
	}

	if updatePlayer.Username != player.Username {
		updatePlayer.Username = player.Username
	}
	if updatePlayer.Role != player.Role {
		updatePlayer.Role = player.Role
	}
	if updatePlayer.Spy != player.Spy {
		updatePlayer.Spy = player.Spy
	}

	return collection.Update(bson.M{"players.playerid": player.PlayerID}, bson.M{"$set": bson.M{"players.$": updatePlayer}})
}

//DelPlayer removes a player matching a suppiled pid
func DelPlayer(pid string) error {
	return collection.Update(bson.M{"players.playerid": bson.ObjectIdHex(pid)}, bson.M{"$pull": bson.M{"players": bson.M{"playerid": bson.ObjectIdHex(pid)}}})
}

//GetPlayers returns an array of players in a game
func GetPlayers(gamecode string) ([]Player, error) {
	game, gdErr := GetGame(gamecode)
	if gdErr != nil {
		return nil, gdErr
	}
	return game.Players, nil
}

//AddGame creates a new game in the DB with the provided gamecode
//Returns error if game already exists or cannot be created
//NOTE: LOCATION NOT SET HERE
func AddGame(gamecode string) error {
	return insertGame(&Game{
		ID:       bson.NewObjectId(),
		GameCode: gamecode,
		Location: "null",
		Players:  nil,
		Active:   false,
	})
}

func SetGame(game *Game) error {
	updateGame, err := GetGame(game.GameCode)
	if err != nil {
		return err
	}

	if updateGame.Location != game.Location {
		err := collection.Update(bson.M{"gamecode": game.GameCode}, bson.M{"$set": bson.M{"location": game.Location}})
		if err != nil {
			return err
		}
	}

	if updateGame.Active != game.Active {
		err := collection.Update(bson.M{"gamecode": game.GameCode}, bson.M{"$set": bson.M{"active": game.Active}})
		if err != nil {
			return err
		}
	}

	if !reflect.DeepEqual(updateGame.Players, game.Players) {
		err := collection.Update(bson.M{"gamecode": game.GameCode}, bson.M{"$set": bson.M{"players": game.Players}})
		if err != nil {
			return err
		}
	}
	return nil
}

//DelGame deletes a game matching the supplied gamecode
func DelGame(gamecode string) error {
	game, err := GetGame(gamecode)
	if err != nil {
		return err
	}

	return collection.Remove(game)
}

//GetGameData returns the Game struct that matches the provided gamecode
//Returns error if game does not exist
func GetGame(gamecode string) (Game, error) {
	gErr := checkGame(gamecode)
	result := Game{}
	if gErr != nil {
		return result, gErr
	}
	err := collection.Find(bson.M{"gamecode": gamecode}).One(&result)
	if err != nil {
		return result, err
	}
	return result, err
}

//GetGameCode returns a gamecode string associated with a supplied pid
func GetGameCode(pid string) (string, error) {
	game := Game{}
	err := collection.Find(bson.M{"players.playerid": bson.ObjectIdHex(pid)}).One(&game)
	return game.GameCode, err
}

//Can be removed if set game works
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
	/*count, err := collection.Find(bson.M{"gamecode": entry.GameCode}).Limit(1).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("GAME_ALREADY_EXISTS")
	}*/
	err := checkGame(entry.GameCode)
	if err != nil {
		return collection.Insert(entry)
	}
	return fmt.Errorf("GAME_EXISTS")
}

//This is the old AddPlayer function. Keeping it just in case
/*func AddPlayer(gamecode, username string) (bson.ObjectId, error) {
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
}*/
