package controller

import (
	"encoding/json"
	"gochat/model"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var UpGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendTo(target model.Sender, msgtype int, msg []byte) error {
	return target.SendMsg(msgtype, msg)
}

func SendToJson(target model.Sender, msg interface{}) error {

	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return SendTo(target, websocket.TextMessage, bytes)
}

func ParseAuthInfo(msg []byte) (info *model.AuthInfo, err error) {
	// error: 此处犯了一个小错误，结构体在没有初始化的时候就拿来使用，这里会报错，结构体在没有初始化的时候不能使用
	info = new(model.AuthInfo)
	err = json.Unmarshal(msg, info)
	return
}

func Auth(info *model.AuthInfo) bool {
	// 是否要从数据库读取info？
	if token, err := model.GetToken(info.Token); err == nil && token.Expires.Sub(time.Now()) > 0 {
		return true
	}
	return false
}

func ParseUserRequest(message []byte) (send *model.SendRequest, err error) {
	send = new(model.SendRequest)
	err = json.Unmarshal(message, &send)
	send.Time = time.Now().Unix()
	return
}
