package config

var MaxUserNum = 30
var MaxChatRoom = 10
var MaxTokenTime = 300 // minute
var ServerPort = ":5000"
var MongoHost = "mongodb://localhost:27017"
var MongoDatabase = "GoChat"
var Mode = 0

const (
	DevMode = iota
	ReleaseMode
)