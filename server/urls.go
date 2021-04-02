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

const Test = true

func Register(router *gin.Engine) {
	router.GET("/ping", pong)
	router.Use(GenerateReqId, AllowOrigin)
	// 管理员api
	adminModule := router.Group("/admin", CheckLoginStatus, Authentication)
	adminModule.POST("/user_list", admin.UserListServer)                           // 获取用户列表
	adminModule.POST("/user_list/csv", admin.UserListCsvServer)                    // 用户列表导出csv
	adminModule.POST("/user_list/print", admin.UserListPrintServer)                // 全部用户列表
	adminModule.POST("/authority/modify")                                          // 管理员修改用户信息
	adminModule.POST("/user_detail", admin.UserDetailServer)                       // 获取用户详细信息（包括权限等）
	adminModule.POST("/user_register/line_chart", admin.GetUserRegisterInfoServer) // 注册人数统计
	// 用户模块
	userModule := router.Group("/user")
	userModule.POST("/register", userAccount.UserRegisterServer)                  // 注册
	userModule.POST("/login/password", userAccount.UserPasswordLoginServer)       // 密码登录
	userModule.POST("/login/verify", userAccount.UserVerifyLoginServer)           // 验证码登录
	userModule.POST("/logout", userAccount.LogoutServer)                          // 注销
	userModule.POST("/modify", CheckLoginStatus, userAccount.ModifyServer)        // 修改用户信息
	userModule.POST("/delete", CheckLoginStatus, userAccount.DeleteServer)        // 删除
	userModule.POST("/apply_role", CheckLoginStatus, userAccount.ApplyRoleServer) // 申请权限
	userModule.POST("/set_avatar", CheckLoginStatus, userAccount.SetAvatarServer) // 设置用户头像
	userModule.POST("/get_profile", CheckLoginStatus, userAccount.GetUserProfile) // 获取用户信息
	userModule.POST("/get_avatar", CheckLoginStatus, userAccount.GetAvatarServer)  // 获取用户头像
	//userModule.POST("/get_user_list", CheckLoginStatus)
	// 工具模块
	toolModule := router.Group("/tool")
	toolModule.POST("/no_auth/generate_verification_code", tool.GenerateVerificationCode) // 发送验证码

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
	if Test {
		ctx.Next()
		return
	}
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
	if Test {
		ctx.Next()
		return
	}
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
	ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, Email")
	ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	method := ctx.Request.Method
	if method == "OPTIONS" {
		ctx.AbortWithStatus(http.StatusNoContent)
	}
	ctx.Next()
}
