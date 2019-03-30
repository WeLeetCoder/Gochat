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
	Id    string
	Name  string
	conn  *websocket.Conn
	Group *Group
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

func (u *User) setGroup(group *Group) {
	u.Group = group
}

func NewUser(id, name string, conn *websocket.Conn) (user *User) {
	user = &User{Id: id, Name: name, conn: nil}
	user.Connect(conn)
	return
}

type Groups map[string]*Group

var UserGroupTable = make(Groups)

func (gs *Groups) Group(name string) *Group {
	return (*gs)[name]
}

func (gs *Groups) HasGroup(name string) (ok bool) {
	_, ok = (*gs)[name]
	return
}

func (gs *Groups) Disconnect(roomname string, user *User) error {
	if gs.HasGroup(roomname) {
		return (*gs)[roomname].Users.Disconnect(user)
	}
	return fmt.Errorf("[Groups disconnect] room %s named error", roomname)
}

func (gs *Groups) Add(group *Group) error {
	if gs.HasGroup(group.RoomName) {
		return fmt.Errorf("room %s is exist. ", group.RoomName)
	}
	(*gs)[group.RoomName] = group
	return nil
}

func (gs *Groups) JoinGroup(name string, user *User) error {
	// 判断用户组是否存在，存在则直接将用户加入进去
	if gs.HasGroup(name) {
		err := (*gs)[name].Add(user)
		if err == nil {
			user.setGroup((*gs)[name])
			return nil
		}
		return err
	}
	// 用户组不存在，则创建新的组，将用户加入进组，并将组附加到组表
	newGroup := NewGroup(name)
	newGroup.Add(user)
	user.setGroup(newGroup)

	// 此处可以不用判断错误，因为上面已经判断过如果有组则加入
	gs.Add(newGroup)
	return nil
}

type Group struct {
	RoomName string
	Users    *Users
}

func (g *Group) Add(user *User) error {
	// 向组加入用户
	return g.Users.Add(user)
}

func (g *Group) SendMsg(msg []byte, msgType int) error {
	return g.Users.SendMsg(msg, msgType)
}

func NewGroup(name string) *Group {
	return &Group{
		RoomName: name,
		Users:    &Users{},
	}
}

func UserConnect(user *User, online bool) error {
	code := 998
	if online {
		code = 999
	}
	return user.Group.SendMsg([]byte(fmt.Sprintf(`{"code":%d, "user": "%s"}`, code, user.Name)), websocket.TextMessage)
}
