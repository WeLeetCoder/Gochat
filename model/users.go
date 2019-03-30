package model

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Users map[string]*User

func (u *Users) SendMsg(msg []byte, msgType int) error {
	for username := range *u {
		if err := (*u)[username].SendMsg(msg, msgType); err != nil {
			return err
		}
	}
	return nil
}
func (u *Users) Disconnect(user *User) error {
	if _, ok := (*u)[user.Id]; ok {
		delete(*u, user.Id)
		return nil
	}
	return fmt.Errorf("user %d is no exist. ", user.Id)
}

func (u *Users) Add(user *User) error {
	if _, ok := (*u)[user.Id]; ok {
		return fmt.Errorf("用户已经存在了")
	}
	(*u)[user.Id] = user
	return nil
}

type User struct {
	Id   string
	Name string
	conn *websocket.Conn
}

// 通过 Connect 函数连接，验证通过之后创建用户，连接用户
func (u *User) Connect(conn *websocket.Conn) error {
	if u.conn != nil {
		return fmt.Errorf("无法连接，用户已经连接了")
	}
	u.conn = conn
	return nil
}

func (u *User) SendMsg(msg []byte, msgType int) error {
	return u.conn.WriteMessage(msgType, msg)
}

func NewUser(id, name string) *User {
	return &User{id, name, nil}
}
