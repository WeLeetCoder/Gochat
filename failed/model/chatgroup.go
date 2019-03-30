package model

import (
	"fmt"
	"time"
)

/*
	聊天组,
	发送消息
	修改名称
	转交群主
*/

type ChatGroup struct {
	Id         string
	Name       string
	CreateTime time.Time
	Owner      *User
	UserLists  *UserLists
}

func (cg *ChatGroup) SendMsg(msgType int, msg []byte) error {
	fmt.Println("send msg to group", cg.Id)
	cg.UserLists.SendMsg(msgType, msg)
	return nil
}

func CreateGroup(name string, owner *User) *ChatGroup {
	id := NewId(name, owner.Name)
	userLists := new(UserLists)
	userLists.Add(owner)
	return &ChatGroup{
		Id:         id,
		Name:       name,
		Owner:      owner,
		CreateTime: time.Now(),
		UserLists:  userLists,
	}
}

func JoinGroup(name string) *ChatGroup {
	id := NewId(name, "")
	userLists := new(UserLists)
	return &ChatGroup{
		Id:         id,
		Name:       name,
		Owner:      nil,
		CreateTime: time.Now(),
		UserLists:  userLists,
	}
}
