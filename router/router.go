package router

import (
	"fmt"
	"gochat/model"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

var UpGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var UserTable = make(model.Users)

func Chat(c *gin.Context) {
	token := c.Param("token")
	fmt.Println(token)
	conn, err := UpGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("[Conn]: ", err.Error())
		return
	}
	defer conn.Close()
	_, msg, err := conn.ReadMessage()
	session, err := model.ParseSession(msg)
	if err != nil {
		return
	}
	user := model.NewUser(session.Id, session.Username)
	err = user.Connect(conn)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = UserTable.Add(user)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = UserTable.SendMsg([]byte(fmt.Sprintf(`{"code": "999", "user": "%s"}`, user.Name)), websocket.TextMessage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			if err = UserTable.Disconnect(user); err != nil {
				fmt.Println("用户下线失败：", err.Error())
			}
			err = UserTable.SendMsg([]byte("用户离线提醒！"), websocket.TextMessage)
			log.Printf("用户 %s Id: %s 下线", user.Name, user.Id)
			break
		}
		message, err := model.ParseMsg(msg)
		if err != nil {
			continue
		}
		message.Sender = user.Id
		message.SenderName = user.Name
		msg, err = message.ToJson()
		if err != nil {
			log.Println("ToJson Error")
			continue
		}
		UserTable.SendMsg(msg, mt)
	}
}

func Join(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	username := c.PostForm("username")
	roomname := c.PostForm("roomname")
	if username == "" || roomname == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 1000,
			"data": "",
			"msg":  "username or roomname empty.",
		})
		return
	}
	fmt.Println(username)
	id := model.NewId(username, roomname)
	authInfo := model.NewInfo(id, username, roomname)
	err := model.SetToken(authInfo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1000,
			"data": "",
			"msg":  "set info failed.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 1000,
		"data": gin.H{
			"Id":       authInfo.Id,
			"Username": authInfo.Username,
			"Roomname": authInfo.Roomname,
			"Token":    authInfo.Token,
		},
		"msg": "",
	})
}

func Users(c *gin.Context) {
	roomName := c.PostForm("rname")
	fmt.Println(roomName, UserTable)
	c.JSON(http.StatusOK, gin.H{
		"code": 1000,
		"data": UserTable,
		"msg":  "",
	})
}

func Rooms(c *gin.Context) {

}

func Broadcast(c *gin.Context) {
	msg := c.PostForm("msg")
	fmt.Println(msg)
}
