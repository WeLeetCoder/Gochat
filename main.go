package main

import (
	"gochat/config"
	"gochat/router"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.DisableConsoleColor()
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.GET("/chat/:token", router.Chat)
	r.POST("/join", router.Join)
	r.POST("/users", router.Users)
	r.POST("/rooms", router.Rooms)
	r.POST("/broadcast", router.Broadcast)
	log.Fatal(http.ListenAndServe(config.ServerPort, r))
}
