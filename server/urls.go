package server

import (
	server "fsbm/server/user_account"
	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine) {
	router.GET("/ping", pong)
	// 用户模块
	userModule := router.Group("/user")
	userModule.POST("/register", server.RegisterServer)
	userModule.POST("/login", server.LoginServer)
	userModule.POST("/logout")
	userModule.POST("/modify")
	userModule.POST("/delete")
}
