package model

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

/*
用户结构体
*/

type User struct {
	Id          string
	Name        string
	Picture     string
	CreateTime  time.Time
	conn        *websocket.Conn
	FriendLists *UserLists
	GroupLists  *GroupLists
}

func (user *User) Connect(conn *websocket.Conn) error {
	if conn == nil {
		return fmt.Errorf("conn is <nil> . ")
	}
	if user.conn != nil {
		return fmt.Errorf("conn already connected. ")
	}
	user.conn = conn
	return nil
}

func (user User) JoinGroup(group *ChatGroup) error {
	return user.GroupLists.Add(group)
}

func (user User) SendMsg(msgType int, msg []byte) error {
	return user.conn.WriteMessage(msgType, msg)
}

func NewId(name, roomName string) string {
	ctx := md5.New()
	ctx.Write([]byte(name + roomName))
	return fmt.Sprintf("%x", ctx.Sum(nil))
}

func RandNewId(name string) string {
	nowTime := strconv.Itoa(int(time.Now().Unix()))
	ctx := md5.New()
	ctx.Write([]byte(name + nowTime))
	return fmt.Sprintf("%x", ctx.Sum(nil))
}

func CreateUser(name string) *User {
	id := RandNewId(name)
	user := new(User)
	user.Id = id
	user.Name = name
	user.CreateTime = time.Now()

	return user
}

func CreateTempUser(name, roomName string) *User {
	//id := NewId(name, roomName)
	//user := new(User)
	//user.Id = id
	//user.Name = name
	//user.CreateTime = time.Now()
	//user.FriendLists = nil
	//user.GroupLists = nil

	id := NewId(name, roomName)
	return &User{
		Id:         id,
		Name:       name,
		CreateTime: time.Now(),
	}
}
