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
			if err = model.UserGroupTable.Disconnect(info.Roomname, user); err != nil {
				fmt.Println("用户下线失败：", err.Error())
			}
			err = model.UserConnect(user, false)
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
		user.Group.SendMsg(msg, mt)
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
	roomName := c.PostForm("rname")
	fmt.Println(roomName, model.UserGroupTable)
	c.JSON(http.StatusOK, gin.H{
		"code": 1000,
		"data": model.UserGroupTable,
		"msg":  "",
	})
}

func Rooms(c *gin.Context) {

}

func Broadcast(c *gin.Context) {
	msg := c.PostForm("msg")
	fmt.Println(msg)
}
