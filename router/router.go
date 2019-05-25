package router

import (
	"Gochat/config"
	"Gochat/model"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

var UpGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Chat(c *gin.Context) {
	c.Status(417)
	token := c.Param("token")
	info, err := model.AuthToken(token)
	if err != nil {
		return
	}

	conn, err := UpGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("[Conn]: ", err.Error())
		return
	}
	defer conn.Close()

	// 新建用户，使用id，和用户名，上面通过token验证，即可
	user := model.NewUser(info.Id, info.Username, conn)
	if err != nil {
		return
	}

	// 将用户加入相应的组
	err = model.UserGroupTable.JoinGroup(info.Roomname, user)
	if err != nil {
		log.Println(err)
		return
	}

	// 提醒组内用户用户上线
	err = model.UserConnect(user, true)
	if err != nil {
		log.Println("UserConnect error.")
		return
	}

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {

			// 用户断开了连接，把用户从用户表中断开连接
			if err = model.UserGroupTable.Disconnect(info.Roomname, user); err != nil {
				fmt.Println("用户下线失败：", err.Error())
			}

			// 提醒组内用户下线
			err = model.UserConnect(user, false)
			log.Printf("用户 %s Id: %s 下线", user.Name, user.Id)
			break
		}

		// 解析用户发送的数据，
		message, err := model.ParseMsg(msg)
		if err != nil {
			continue
		}

		msg, err = message.ToJson(user)
		if err != nil {
			log.Println("ToJson Error")
			continue
		}

		// 在用户组中发送消息
		user.SendGroup(msg, mt)
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
	c.Header("Access-Control-Allow-Origin", "*")
	roomName := c.PostForm("rname")
	token := c.PostForm("token")

	fmt.Println(roomName, token)
	bytes, err := model.UserGroupTable.Group(roomName)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1006,
			"data": bytes,
			"msg":  err,
		})
		return
	}
	bytes.GenerateUList()
	c.JSON(http.StatusOK, gin.H{
		"code": 1000,
		"data": bytes,
		"msg":  err,
	})
}

func Rooms(c *gin.Context) {
	fmt.Println(model.UserGroupTable)
	c.JSON(http.StatusOK, gin.H{
		"code": 1000,
		"data": model.UserGroupTable,
		"msg":  "this is group table.",
	})
}

func Broadcast(c *gin.Context) {
	token := c.PostForm("token")
	_, err := model.AuthToken(token)
	if err != nil {
		return
	}
	msg := c.PostForm("msg")
	newMsg, err := model.NewMsg(config.BroadcastName, "System", msg, time.Now())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1006,
			"data": "",
			"msg":  "Broadcast failed.",
		})
		return
	}

	for _, user := range model.UserGroupTable {
		user.SendMsg(newMsg, websocket.TextMessage)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1000,
		"data": "",
		"msg":  "Broadcast successful.",
	})
}