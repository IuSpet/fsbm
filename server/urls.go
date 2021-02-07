package server

import (
	userAccount "fsbm/server/user_account"
	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine) {
	router.GET("/ping", pong)
	// 用户模块
	userModule := router.Group("/user")
	userModule.POST("/register", userAccount.UserRegisterServer)
	userModule.POST("/login", userAccount.UserLoginServer)
	userModule.POST("/logout")
	userModule.POST("/modify")
	userModule.POST("/delete")
}
