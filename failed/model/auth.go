package model

import "time"

type AuthInfo struct {
	Username string
	RoomName string
	Token    string
}

type Token struct {
	Username, RoomName, Token, Id string
	Expires                       time.Time
}
