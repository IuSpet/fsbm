package server

import (
	"fmt"
	"fsbm/server/tool"
	userAccount "fsbm/server/user_account"
	"fsbm/util"
	"fsbm/util/redis"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func Register(router *gin.Engine) {
	router.GET("/ping", pong)
	// 用户模块
	userModule := router.Group("/user", GenerateReqId)
	userModule.POST("/register", userAccount.UserRegisterServer)
	userModule.POST("/login/password", userAccount.UserPasswordLoginServer)
	userModule.POST("/login/verify", userAccount.UserVerifyLoginServer)
	userModule.POST("/logout", userAccount.LogoutServer)
	userModule.POST("/modify", CheckLoginStatus, userAccount.ModifyServer)
	userModule.POST("/delete", CheckLoginStatus, userAccount.DeleteServer)
	// 工具模块
	toolModule := router.Group("/tool")
	toolModule.POST("/no_auth/generate_verification_code", tool.GenerateVerificationCode)
}

func GenerateReqId(ctx *gin.Context) {
	// 用请求到达的时间戳当req_id
	now := time.Now()
	ctx.Set("req_id", strconv.FormatInt(now.Unix(), 10))
}

func CheckLoginStatus(ctx *gin.Context) {
	email := ctx.GetHeader("email")
	key := fmt.Sprintf(util.UserLoginTemplate, email)
	res, err := redis.GetWithRetry(ctx, key)
	if err != nil || res == "" {
		ctx.Abort()
		util.ErrorJson(ctx, util.UserNotLogin, "用户未登陆")
		return
	}
}
