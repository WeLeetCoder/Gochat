package model

import (
	"encoding/json"
	"fmt"
	"html"
	"strings"
	"time"
)

// 定义从websocket中接受的数据，并解析

type Message struct {
	SenderName string
	Sender     string
	Content    string
	Time       string
}

func (m Message) ToJson(user *User) ([]byte, error) {
	m.SenderName = user.Name
	m.Sender = user.Id
	return json.Marshal(m)
}

func NewMsg(name, sender, content string, time time.Time) ([]byte, error) {
	msg := new(Message)
	msg.Time = time.Format("2006/01/02 15:04:05")
	msg.Sender = name
	msg.Content = content
	msg.SenderName = sender
	return json.Marshal(msg)
}

type Session struct {
	Username string
	Id       string
}

func ParseSession(msg []byte) (s *Session, err error) {
	s = new(Session)
	err = json.Unmarshal(msg, s)
	return
}

func ParseMsg(m []byte) (msg *Message, err error) {
	msg = new(Message)
	msg.Time = time.Now().Format("2006/01/02 15:04:05")
	err = json.Unmarshal(m, msg)
	content := html.EscapeString(msg.Content)
	msg.Content = strings.Replace(fmt.Sprintf("<pre>%s</pre>", content), "\n", "<br>", -1)
	return
}
