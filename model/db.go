package model

import (
	"gochat/config"
	"log"
	"os"

	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"
)

var MongoDB *mgo.Database

func initMongo() {
	session, err := mgo.Dial(config.MongoHost)
	if err != nil {
		log.Println("[MongoDB Init]", err.Error())
		os.Exit(-1)
	}
	session.SetMode(mgo.Monotonic, true)
	MongoDB = session.DB(config.MongoDatabase)
}

func SetToken(info *AuthInfo) error {
	collection := MongoDB.C(config.MongoTokenTab)
	return collection.Insert(info)
}

func GetToken(token string) (info *AuthInfo, err error) {
	collection := MongoDB.C(config.MongoTokenTab)
	info = new(AuthInfo)
	err = collection.Find(bson.M{"token": token}).One(info)
	return
}

func init() {
	initMongo()
}
