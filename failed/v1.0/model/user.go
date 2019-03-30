package model

import "github.com/gorilla/websocket"

type User struct {
	Username string
	Conn     *websocket.Conn
}

func CreateUser(name string) *User {
	return &User{
		Username: name,
	}
}

func (user *User) Connect(conn *websocket.Conn) {
	user.Conn = conn
}

func (user *User) IsConnect() bool {
	if user.Conn == nil {
		return false
	}
	return true
}

type Join struct {
	Rname    string
	Username string
}
