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

//func SetToken(info *AuthInfo) (token *Token, err error) {
//	tokenTable := MongoDB.C("token")
//	tokenStr := RandNewId(time.Now().String())
//	token = &Token{
//		Id:       NewId(info.Username, info.RoomName),
//		Username: info.Username,
//		RoomName: info.RoomName,
//		Expires:  time.Now().Add(30 * time.Minute),
//	}
//	token.Token = tokenStr
//	err = tokenTable.Insert(token)
//	return
//}
//
//func GetToken(tokenStr string) (token *Token, err error) {
//	tokenTable := MongoDB.C("token")
//	err = tokenTable.Find(bson.M{"token": tokenStr}).One(&token)
//	return
//}

func init() {
	initMongo()
}
