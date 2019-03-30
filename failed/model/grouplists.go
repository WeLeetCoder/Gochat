package model

import "fmt"

/*
群组列表，包含了所有的群组
通过 map 实现一个列表，键值是聊天组的id值，用id值作为键，string 是用户名称，即可以实现两个组名相同，不知此处是否有必要，两个组名相同
增加：判断群组是否存在，如果存在则，增加失败，两个群名称不能相同
删除：从当前的群组列表删除，而不是从服务器上删除，判断是否存在，存在则删除，不存在则返回一个错误
查询：判断是否存在，存在则返回一个群组，不存在则群组返回空，并返回一个err

*/

type GroupLists map[string]*ChatGroup

func (g *GroupLists) Add(group *ChatGroup) error {
	fmt.Println(g, group)
	if g.HasMember(group.Id) {
		return fmt.Errorf("add group name %s failed. ", group.Name)
	}
	(*g)[group.Id] = group
	return nil
}

func (g GroupLists) Delete(id string) error {
	if g.HasMember(id) {
		delete(g, id)
		return nil
	}
	return fmt.Errorf("delete group id %s failed. ", id)
}

func (g GroupLists) Update(id string, name string) error {
	if g.HasMember(id) {
		g[id].Name = name
		return nil
	}
	return fmt.Errorf("update group id %s failed. ", id)
}

func (g GroupLists) HasMember(id string) (ok bool) {
	_, ok = g[id]
	return
}

func (g GroupLists) GetMember(id string) (*ChatGroup, error) {
	if g.HasMember(id) {
		return g[id], nil
	}
	return nil, fmt.Errorf("get group member %s faild. ", id)
}

func (g GroupLists) GetMembers() []string {
	users := make([]string, 0)
	for name := range g {
		users = append(users, name)
	}
	return users
}
