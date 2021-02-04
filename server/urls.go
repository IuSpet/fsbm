package server

import "github.com/gin-gonic/gin"

func Register(router *gin.Engine) {
	router.GET("/ping", pong)
	// 用户模块
	userModule := router.Group("/user")
	userModule.POST("/register")
	userModule.POST("/login")
	userModule.POST("/logout")
	userModule.POST("/modify")
	userModule.POST("/delete")
}
