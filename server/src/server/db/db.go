package db

import (
	"gopkg.in/mgo.v2"
)

type DBO struct {
	Server         string
	Database       string
	GameCollection string
}

var db *mgo.Database

func Connect(dbo DBO) error {
	session, err := mgo.Dial(dbo.Server)
	if err != nil {
		return err
	}

	db = session.DB(dbo.Database)
	defer session.Close()

	return nil
}
