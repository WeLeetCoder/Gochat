package config

const MaxUserNum = 60
const MaxChatRoom = 10
const MaxTokenTime = 300 // minute
const ServerPort = ":5000"
const MongoHost = "mongodb://localhost:27017"
const MongoDatabase = "GoChat"
const MongoTokenTab = "token"
const Mode = 0
const TokenSecureCode = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-~|<>?/';:"
const BroadcastName = "GoChat"

const (
	DevMode = iota
	ReleaseMode
)
