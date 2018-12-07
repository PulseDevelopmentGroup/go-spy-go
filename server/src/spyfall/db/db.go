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
	Spy      bool          `bson:"spy" json:"spy"`
}

type Game struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	GameCode string        `bson:"gamecode" json:"gamecode"`
	Location string        `bson:"location" json:"location"`
	Players  []Player      `bson:"players" json:"players"`
}

var collection *mgo.Collection

func Connect(dbo *DBO) error {
	session, err := mgo.Dial(dbo.Server)
	if err != nil {
		return err
	}

	db := session.DB(dbo.Database)
	collection = db.C(dbo.GameCollection)

	return nil
}

func AddPlayer(gamecode, username string) (string, error) {
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
	} else {
		pid := bson.NewObjectId()
		err := collection.Update(bson.M{"gamecode": gamecode}, bson.M{"$push": bson.M{"players": &Player{
			PlayerID: pid,
			Username: username,
			Spy:      false,
		}}})
		return pid.Hex(), err
	}
}

func GetPlayers(gamecode string) ([]string, error) {
	game, gdErr := GetGameData(gamecode)
	players := []string{}
	if gdErr != nil {
		if gdErr == fmt.Errorf("NO_GAME_EXISTS") {
			return nil, fmt.Errorf("NO_GAME_EXISTS")
		}
		return nil, gdErr
	}
	fmt.Println(len(game.Players))
	for i := 0; i < len(game.Players); i++ {
		players = append(players, game.Players[i].Username)
	}

	return players, nil
}

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

func NewGame(gamecode, location string) error {
	return insertEntry(&Game{
		ID:       bson.NewObjectId(),
		GameCode: gamecode,
		Location: location,
		Players:  nil,
	})
}

//TODO Test this more thouroughly    < Also, definitely spelled that wrong
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

func insertEntry(entry *Game) error {
	count, err := collection.Find(bson.M{"gamecode": entry.GameCode}).Limit(1).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("GAME_ALREADY_EXISTS")
	}

	return collection.Insert(entry)
}
