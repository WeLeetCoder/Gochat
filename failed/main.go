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
	r := gin.Default()
	//r.GET("/", router.Home)
	r.GET("/grouplist", router.ListGroups)
	r.GET("/userlist", router.ListUsers)
	r.POST("/join", router.Join)
	if config.Mode == config.DevMode {
		r.GET("/chat", router.Chat)
		r.GET("/verify", router.GetToken)
	}
	log.Fatal(http.ListenAndServe(config.ServerPort, r))
}
