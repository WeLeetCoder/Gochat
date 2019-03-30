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
	Time       int64
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
	msg.Time = time.Now().Unix()
	err = json.Unmarshal(m, msg)
	return
}
