package model

import (
	"gochat/config"
	"log"
	"os"
	"time"

	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"
)

/*
	用户和群组的主要存储地，通过id来获取用户，或者群组
*/

/*
连接数据库，获取数据库中的用户信息，数据库中存储所有的用户，还有用户组
用户已经上线判断？
用户是否在线？
此处是模拟用户上线
*/

// 数据库表该如何创建？
// 该创建哪些表？
// 用户表                                 临时用户？是否需要进行创建？ 临时的用户组如何维护

// 此处遇到了巨大的坑，忘记了userlist 是map，用new来分配就gg了，报了n次错误
// var UserTable = new(UserLists)
// var GroupTable = new(GroupLists)

var UserTable UserLists = make(UserLists)
var GroupTable GroupLists = make(GroupLists)

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

func SetToken(info *AuthInfo) (token *Token, err error) {
	tokenTable := MongoDB.C("token")
	tokenStr := RandNewId(time.Now().String())
	token = &Token{
		Id:       NewId(info.Username, info.RoomName),
		Username: info.Username,
		RoomName: info.RoomName,
		Expires:  time.Now().Add(30 * time.Minute),
	}
	token.Token = tokenStr
	err = tokenTable.Insert(token)
	return
}

func GetToken(tokenStr string) (token *Token, err error) {
	tokenTable := MongoDB.C("token")
	err = tokenTable.Find(bson.M{"token": tokenStr}).One(&token)
	return
}

func init() {
	initMongo()
}
