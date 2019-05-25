package config

const (
	MaxUserNum      = 60
	MaxChatRoom     = 10
	MaxTokenTime    = 300 // minute
	ServerPort      = ":5000"
	MongoHost       = "mongodb://localhost:27017"
	MongoDatabase   = "GoChat"
	MongoTokenTab   = "token"
	Mode            = 0
	TokenSecureCode = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-~|<>?/';:"
	BroadcastName   = "GoChat"
)

const (
	DevMode = iota
	ReleaseMode
)
