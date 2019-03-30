package model

import (
	"fmt"
	"time"
)

// 此处需要注意，之前用的是map[string]*Chatroom 报错了，然后用make分配了内存，正确
var Rlist RoomList = make(map[string]*Chatroom)

type RoomList map[string]*Chatroom

// 删除整个房间
func (rl RoomList) Delete(name string) error {
	if _, ok := rl[name]; ok {
		delete(rl, name)
		return nil
	}
	return fmt.Errorf("room %s is not exist.", name)
}

func (rl RoomList) Room(name string) (*Chatroom, error) {
	room, ok := rl[name]
	if !ok {
		return nil, fmt.Errorf("room %s is not exist!", name)
	}
	return room, nil
}

// 添加房间，必须要有房主
func (rl RoomList) Add(name string, owner *User) *Chatroom {
	room := NewRoom(name, owner)
	rl[name] = room
	return room
}

// 判断房间是否存在
func (rl RoomList) IsExist(name string) bool {
	_, ok := rl[name]
	return ok
}

// 列出所有房间的名字
func (rl RoomList) RoomList() (roomlist []string) {
	for name := range rl {
		roomlist = append(roomlist, name)
	}
	return roomlist
}

// 聊天室结构体
type Chatroom struct {
	RoomId     int64
	RoomName   string
	CreateTime time.Time
	Owner      *User
	Users      Clients // 存储了每个用户，键名为用户名
}

// 为房间添加用户，用户不可以重复
func (room *Chatroom) Join(user *User) (*Chatroom, error) {
	if room.IsExist(user.Username) {
		return nil, fmt.Errorf("user %s add failure", user.Username)
	}
	room.Users[user.Username] = user
	return room, nil
}

func (room *Chatroom) DeleteUser(user string) error {
	err := room.Users.Delete(user)
	room.destroy()
	return err
}

func (room *Chatroom) IsExist(user string) bool {
	_, ok := room.Users[user]
	fmt.Println("---->[exist]", ok)
	return ok
}

func (room *Chatroom) Broadcast(mt int, msg []byte) error {
	// 广播时尝试删除房间？
	if err := room.destroy(); err != nil {
		return err
	}
	return broadcast(mt, msg, room)
}

func (room Chatroom) destroy() error {
	if len(room.Users) == 0 { // 如果用户数量为 0 的话，销毁这个房间
		if err := Rlist.Delete(room.RoomName); err != nil { // 删除这个房间
			return err
		}
	}
	return nil
}

type Clients map[string]*User

func (clients Clients) hasUser(name string) bool {
	if _, ok := clients[name]; ok {
		return true
	}
	return false
}

// 列出房间内的用户
func (clients Clients) ListClients() (users []string) {
	for name := range clients {
		users = append(users, name)
	}
	return
}

func (clients Clients) Delete(name string) error {

	if clients.hasUser(name) {
		delete(clients, name)
		return nil
	}
	return fmt.Errorf("user name is not exist")
}

func (clients Clients) GetUser(name string) (*User, error) {
	if user, ok := clients[name]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("user %s is not exist.", name)
}

// 创建新房间
func NewRoom(name string, owner *User) *Chatroom {
	var id = time.Now().Unix()
	room := &Chatroom{
		RoomId:     id,
		RoomName:   name,
		CreateTime: time.Now(),
		Owner:      owner,
		Users:      make(map[string]*User),
	}
	room.Join(owner)
	return room
}

func broadcast(mt int, msg []byte, target *Chatroom) error {
	for _, user := range target.Users {
		if err := user.Conn.WriteMessage(mt, msg); err != nil {
			return err
		}
	}
	return nil
}
