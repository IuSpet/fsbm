package server

import (
	"fmt"
	"fsbm/server/admin"
	"fsbm/server/tool"
	userAccount "fsbm/server/user_account"
	"fsbm/util"
	"fsbm/util/auth"
	"fsbm/util/logs"
	"fsbm/util/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func Register(router *gin.Engine) {
	router.GET("/ping", pong)
	router.Use(GenerateReqId, AllowOrigin)
	// 管理员api
	adminModule := router.Group("/admin", CheckLoginStatus, Authentication)
	adminModule.POST("/user_list", admin.UserListServer)
	adminModule.POST("/authority/modify")
	adminModule.POST("/user_detail", admin.UserDetailServer)
	// 用户模块
	userModule := router.Group("/user")
	userModule.POST("/register", userAccount.UserRegisterServer)
	userModule.POST("/login/password", userAccount.UserPasswordLoginServer)
	userModule.POST("/login/verify", userAccount.UserVerifyLoginServer)
	userModule.POST("/logout", userAccount.LogoutServer)
	userModule.POST("/modify", CheckLoginStatus, userAccount.ModifyServer)
	userModule.POST("/delete", CheckLoginStatus, userAccount.DeleteServer)
	userModule.POST("/apply_role", CheckLoginStatus, userAccount.ApplyRoleServer)
	userModule.POST("/set_avatar", CheckLoginStatus, userAccount.SetAvatarServer)
	userModule.POST("/get_user_list",CheckLoginStatus,)
	// 工具模块
	toolModule := router.Group("/tool")
	toolModule.POST("/no_auth/generate_verification_code", tool.GenerateVerificationCode)

}

// 生成唯一请求ID
func GenerateReqId(ctx *gin.Context) {
	// 用请求到达的时间戳当req_id
	now := time.Now()
	ctx.Set("req_id", strconv.FormatInt(now.UnixNano(), 10))
	ctx.Next()
}

// 检查登录状态
func CheckLoginStatus(ctx *gin.Context) {
	email := ctx.GetHeader("email")
	ctx.Set("email", email)
	key := fmt.Sprintf(util.UserLoginTemplate, email)
	res, err := redis.GetWithRetry(ctx, key)
	if err != nil || res == "" {
		util.ErrorJson(ctx, util.UserNotLogin, "用户未登陆")
		ctx.Abort()
		return
	}
	ctx.Next()
}

// 检查接口权限
func Authentication(ctx *gin.Context) {
	email := ctx.GetString("email")
	userRoleSubject, err := auth.NewUserRoleSubject(email)
	if err != nil {
		logs.CtxError(ctx, "db error: %+v", err)
		util.ErrorJson(ctx, util.DbError, "内部错误")
		ctx.Abort()
		return
	}
	path := ctx.FullPath()
	// api权限检查
	hasPermission := userRoleSubject.HasPermission(ctx, AllPathPermission[path])
	if !hasPermission {
		logs.CtxInfo(ctx, "%s has no permission. permission: %+v", email, AllPathPermission[path])
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	ctx.Next()
}

// 允许跨域调用
func AllowOrigin(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	method := ctx.Request.Method
	if method == "OPTIONS" {
		ctx.AbortWithStatus(http.StatusNoContent)
	}
	ctx.Next()
}
