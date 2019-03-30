package model

import "fmt"

/*
	用户列表也维护一张，一张所有用户的用户列表
	用户列表结构体
	获取用户，删除用户，判断用户
*/

type UserLists map[string]*User

func (u *UserLists) Add(user *User) error {
	if u.HasMember(user.Id) {
		return fmt.Errorf("add user id %s failed. ", user.Id)
	}

	//(*u)[user.Id] = user
	return nil
}

func (u UserLists) Delete(id string) error {
	if u.HasMember(id) {
		delete(u, id)
		return nil
	}
	return fmt.Errorf("delete user id %s failed. ", id)
}

func (u UserLists) Update(id string, name string) error {
	if u.HasMember(id) {
		u[id].Name = name
		return nil
	}
	return fmt.Errorf("update user id %s failed. ", id)
}

func (u UserLists) HasMember(id string) (ok bool) {
	_, ok = u[id]
	return
}

func (u UserLists) GetMember(id string) (*User, error) {
	if u.HasMember(id) {
		return u[id], nil
	}
	return nil, fmt.Errorf("get user %s faild. ", id)
}

func (u UserLists) GetMembers() []string {
	users := make([]string, 0)
	for name := range u {
		users = append(users, name)
	}
	return users
}

func (u *UserLists) SendMsg(msgType int, msg []byte) error {
	for key := range *u {
		(*u)[key].SendMsg(msgType, msg)
	}
	return nil
}
