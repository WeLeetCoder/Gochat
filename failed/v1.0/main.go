package v1_0

import (
	//"gochat/config"
	//"gochat/router"
	//"net/http"

	"gochat/config"
	"gochat/router"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	ginRouter()
}

func ginRouter() {
	r := gin.Default()
	r.Static("/static", "./website")
	r.GET("/chat", router.Chat)            // 这个就是 websocket
	r.POST("/join", router.JoinRoom)       // 使用ajax 新建房间
	r.POST("/roomusers", router.RoomUsers) // 查询房间内人数
	r.POST("/roomlist", router.Rooms)      // 列出所有的房间
	r.POST("/offline", router.Offline)     // 下线某个用户
	s := &http.Server{
		Addr:           config.ServerPort,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
