package router

import (
	"encoding/json"
	"fmt"
	"gochat/model"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

// 首页
func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "website/index.html")
}

var chatroomlist model.RoomList

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func auth(joinInfo model.Join, ws *websocket.Conn) (room *model.Chatroom, user *model.User, err error) {
	room, err = model.Rlist.Room(joinInfo.Rname)
	if err != nil {
		err = ws.WriteMessage(websocket.TextMessage, []byte("join room error1"))
		panic(err)
		return nil, nil, fmt.Errorf("room %s is not exist", joinInfo.Rname)
	}
	user, err = room.Users.GetUser(joinInfo.Username)
	if err != nil {
		err = ws.WriteMessage(websocket.TextMessage, []byte("join room error2"))
		panic(err)
		return nil, nil, fmt.Errorf("user %s is not exist", joinInfo.Username)
	}
	if user.IsConnect() {
		err = ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("join room %s Failed. User %s is loged.", joinInfo.Rname, joinInfo.Username)))
		return nil, nil, fmt.Errorf("user %s is signed", joinInfo.Username)
	}
	return room, user, nil
}

// 会话
func Chat(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	//mt, message, err := ws.ReadMessage()
	mt, message, err := ws.ReadMessage()
	if err != nil {
		fmt.Println("read", err)
		return
	}
	var joinInfo model.Join
	json.Unmarshal(message, &joinInfo)

	room, user, err := auth(joinInfo, ws)
	if err != nil {
		return
	}
	user.Connect(ws)

	//err = ws.WriteMessage(mt, []byte("join room Successful"))
	err = room.Broadcast(mt, message)

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("read", err)
			if err = room.DeleteUser(joinInfo.Username); err != nil {
				log.Println("error", err)
			}
			break
		}

		// 写入 ws 数据
		err = room.Broadcast(mt, message)
		if err != nil {
			fmt.Println("write", err)
			break
		}
	}
}

func JoinRoom(c *gin.Context) {
	name := c.PostForm("username")
	rname := c.PostForm("rname")
	fmt.Printf("user:%T %v room: %T %v", name, name, rname, rname)
	if name == "" || rname == "" {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "username and room name is empty.",
			"code": 1001,
		})
		return
	}
	user := model.CreateUser(name)
	if model.Rlist.IsExist(rname) {
		if room, err := model.Rlist[rname].Join(user); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "Join chatroom Fail(username is exist).",
				"code": 1002,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "create room success",
				"id":   room.RoomId,
				"code": 1000,
			})
		}

		return
	}
	room := model.Rlist.Add(rname, user)
	c.JSON(http.StatusOK, gin.H{
		"msg":  "create room success",
		"id":   room.RoomId,
		"code": 1000,
	})
}

func RoomUsers(c *gin.Context) {
	rname := c.PostForm("rname")
	room, err := model.Rlist.Room(rname)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	userList := room.Users.ListClients()
	fmt.Println(userList)
	c.JSON(http.StatusOK, gin.H{
		"users": userList,
	})
}

func Rooms(c *gin.Context) {
	rList := model.Rlist.RoomList()
	c.JSON(http.StatusOK, gin.H{
		"roomList": rList,
	})
}

func Offline(c *gin.Context) {
	name := c.PostForm("username")
	rname := c.PostForm("rname")

	room, err := model.Rlist.Room(rname)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err = room.DeleteUser(name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":      true,
		"userlist": room.Users.ListClients(),
	})
}
