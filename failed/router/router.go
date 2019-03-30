package router

import (
	"fmt"
	"gochat/controller"
	"gochat/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {

}

func ListGroups(c *gin.Context) {
	mbs := model.GroupTable.GetMembers()
	c.JSON(http.StatusOK, gin.H{
		"mbs": mbs,
	})
}

func ListUsers(c *gin.Context) {
	mbs := model.UserTable.GetMembers()
	c.JSON(http.StatusOK, gin.H{
		"mbs": mbs,
	})
}

func CreateAccount(c *gin.Context) {

}

func Chat(c *gin.Context) {
	conn, err := controller.UpGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	// Todo: websocket 验证用户登录token，与用户名是否匹配，不匹配直接退出
	// Todo: 判断token的有效时间，如果Token超时了发送会话将断开连接
	// 循环接收数据，对接收的数据进行分类处理，接收的数据主要是，发送者？ 接收者？ 发送时间？ token？
	// 如果数据是有效的话，就继续
	_, msg, err := conn.ReadMessage()
	authInfo, err := controller.ParseAuthInfo(msg)
	if err != nil {
		fmt.Println("验证信息解析错误")
		return
	}

	// 判断token是否有效
	tokenVail := controller.Auth(authInfo)
	if !tokenVail {
		fmt.Println("token 无效")
		return
	}

	user := model.CreateTempUser(authInfo.Username, authInfo.RoomName)
	err = user.Connect(conn)
	if err != nil {
		log.Println("user Connect:", err.Error())
		return
	}

	err = model.UserTable.Add(user)
	model.UserTable[user.Id] = user
	if err != nil {
		log.Println("add user error", err.Error())
		return
	}

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			model.UserTable.Delete(user.Id)
			log.Println("read msg error:", err)
			break
		}
		sendRq, err := controller.ParseUserRequest(msg)
		if err != nil {
			log.Println("send request parse error")
			continue
		}
		sender, err := sendRq.GetReceiver()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		err = controller.SendTo(sender, mt, msg)
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func GetToken(c *gin.Context) {
	tokenStr := c.Query("token")
	fmt.Println(tokenStr)
	token, err := model.GetToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  1006,
			"token": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    1000,
		"token":   token.Token,
		"now":     time.Now(),
		"expires": token.Expires,
		"vail":    time.Now().Sub(token.Expires) < 0,
	})
}

func Join(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	username := c.PostForm("username")
	roomname := c.PostForm("roomname")
	// 先判断用户是否已经存在了，如果用户已经存在了将返回错误，告诉客户端已经有该用户了
	// 如果用户没有则返回一个token，token对应用户名，一个用户持有一个token
	if username == "" && roomname == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 1002,
			"data": "",
			"msg":  "username and roomname don't empty. ",
		})
		return
	}
	if model.UserTable.HasMember(model.NewId(username, roomname)) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1001,
			"data": "",
			"msg":  "username is exist.",
		})
		return
	}
	authinfo := &model.AuthInfo{
		Username: username,
		RoomName: roomname,
	}
	token, err := model.SetToken(authinfo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1004,
			"msg":  err.Error(),
			"data": "",
		})
		return
	}
	fmt.Println(token)
	c.JSON(http.StatusOK, gin.H{
		"code": 1000,
		"data": gin.H{
			"Id":       token.Id,
			"Token":    token.Token,
			"Username": token.Username,
			"chatLists": []gin.H{
				gin.H{
					"name": token.RoomName,
				},
			},
		},
		"msg": "",
	})
	return

}
