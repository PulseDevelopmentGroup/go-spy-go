package main

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type dbo struct {
	Server         string
	Database       string
	GameCollection string
}

type player struct {
	PlayerID bson.ObjectId `bson:"playerid" json:"playerid"`
	Username string        `bson:"username" json:"username"`
	Spy      bool          `bson:"spy" json:"spy"`
}

type gameTemplate struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	GameCode string        `bson:"gamecode" json:"gamecode"`
	Location string        `bson:"location" json:"location"`
	Players  []player      `bson:"players" json:"players"`
}

var collection *mgo.Collection

func connect(dbo *dbo) error {
	session, err := mgo.Dial(dbo.Server)
	if err != nil {
		return err
	}

	db := session.DB(dbo.Database)
	collection = db.C(dbo.GameCollection)

	return nil
}

func addPlayer(gamecode, username string) error {
	gameCount, err := collection.Find(bson.M{"gamecode": gamecode}).Limit(1).Count()
	if err != nil {
		return err
	}
	if gameCount == 0 {
		return fmt.Errorf("No game exists with gamecode: %s", gamecode)
	}
	//Neither of the following work so...
	/*usernameCount, err := collection.Find(bson.M{"gamecode": gamecode, "players": bson.M{"username": username}}).Limit(1).Count() //Probably need to get more specific, filtering by gamecode first
	usernameCount, err := collection.Find(bson.M{"gamecode": gamecode}).Limit(1).Select(bson.M{"players": bson.M{"$elemMatch": bson.M{"username": username}}}).Count()
	if err != nil {
		return err
	}*/
	if /*usernameCount */ 0 > 1 {
		return err
	}

	return collection.Update(bson.M{"gamecode": gamecode}, bson.M{"$push": bson.M{"players": &player{ //This works, so that's neat
		PlayerID: bson.NewObjectId(),
		Username: username,
		Spy:      false,
	}}})
}

func newGame(gamecode, location string) error {
	return insertEntry(&gameTemplate{
		ID:       bson.NewObjectId(),
		GameCode: gamecode,
		Location: location,
		Players:  nil,
	})
}

func insertEntry(entry *gameTemplate) error {
	count, err := collection.Find(bson.M{"gamecode": entry.GameCode}).Limit(1).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("Resource already exists with this gamecode: %s", entry.GameCode)
	}

	return collection.Insert(entry)
}
