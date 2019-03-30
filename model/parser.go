package model

import (
	"encoding/json"
	"time"
)

// 定义从websocket中接受的数据，并解析

type Message struct {
	SenderName string
	Sender     string
	Content    string
	Time       string
}

func (m Message) ToJson() ([]byte, error) {
	return json.Marshal(m)
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
	return
}