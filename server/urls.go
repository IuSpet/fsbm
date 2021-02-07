package server

import (
	"fsbm/server/tool"
	userAccount "fsbm/server/user_account"
	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine) {
	router.GET("/ping", pong)
	// 用户模块
	userModule := router.Group("/user")
	userModule.POST("/register", userAccount.UserRegisterServer)
	userModule.POST("/login/password", userAccount.UserPasswordLoginServer)
	userModule.POST("/login/verify", userAccount.UserVerifyLoginServer)
	userModule.POST("/logout")
	userModule.POST("/modify")
	userModule.POST("/delete")
	// 工具模块
	toolModule := router.Group("/tool")
	toolModule.POST("/no_auth/generate_verification_code", tool.GenerateVerificationCode)
}
