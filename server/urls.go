package server

import (
	"fmt"
	"fsbm/server/admin"
	"fsbm/server/alarm"
	"fsbm/server/authority"
	"fsbm/server/dashboard"
	"fsbm/server/shop"
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
	userModule.POST("/register", userAccount.UserRegisterServer)                    // 注册
	userModule.POST("/login/password", userAccount.UserPasswordLoginServer)         // 密码登录
	userModule.POST("/login/verify", userAccount.UserVerifyLoginServer)             // 验证码登录
	userModule.POST("/logout", userAccount.LogoutServer)                            // 注销
	userModule.POST("/modify", CheckLoginStatus, userAccount.ModifyServer)          // 修改用户信息
	userModule.POST("/delete", CheckLoginStatus, userAccount.DeleteServer)          // 删除
	userModule.POST("/apply_role", CheckLoginStatus, userAccount.ApplyRoleServer)   // 申请权限
	userModule.POST("/set_avatar", CheckLoginStatus, userAccount.SetAvatarServer)   // 设置用户头像
	userModule.POST("/get_profile", CheckLoginStatus, userAccount.GetUserProfile)   // 获取用户信息
	userModule.POST("/get_avatar", CheckLoginStatus, userAccount.GetAvatarServer)   // 获取用户头像
	userModule.POST("/get_roles", CheckLoginStatus, userAccount.GetUserRolesServer) // 获取用户角色
	//userModule.POST("/get_info",CheckLoginStatus,userAccount.GetInfoServer)
	//userModule.POST("/get_user_list", CheckLoginStatus)
	// 店铺与设备模块
	shopModule := router.Group("/shop", CheckLoginStatus, Authentication)
	shopModule.POST("/add_shop", shop.AddShopServer)                              // 增加店铺
	shopModule.POST("/shop_info", shop.GetShopInfoServer)                         // 店铺信息
	shopModule.POST("/shop_list", shop.GetShopListServer)                         // 店铺列表
	shopModule.POST("/shop_list/csv", shop.GetShopListCsvServer)                  // 店铺列表csv
	shopModule.POST("/shop_list/print", shop.GetShopListPrintServer)              // 店铺列表pdf
	shopModule.POST("/device/add_monitor", shop.AddMonitorServer)                 // 增加监控
	shopModule.POST("/device/monitor_list", shop.GetMonitorListServer)            // 监控列表
	shopModule.POST("/device/monitor_list/csv", shop.GetMonitorLIstCsvServer)     // 监控列表csv
	shopModule.POST("/device/monitor_list/print", shop.GetMonitorListPrintServer) // 监控列表pdf
	shopModule.POST("/device/live_wall_src", shop.GetLiveWallSrcServer)           // 直播墙源
	shopModule.POST("/shop_list_by_email", shop.GetShopListByEmailServer)         // 某用户负责店铺
	// 权限管理模块
	authModule := router.Group("/auth", CheckLoginStatus, Authentication)
	authModule.POST("/role_list", authority.GetRoleListServer)                     // 系统内所有角色列表
	authModule.POST("/user_role_list", authority.GetUserRoleListServer)            // 用户角色列表
	authModule.POST("/apply_role", authority.ApplyRoleServer)                      // 申请角色
	authModule.POST("/apply_order_list", authority.ApplyRoleListServer)            // 申请角色工单列表
	authModule.POST("/apply_order_list/csv", authority.ApplyRoleListCsvServer)     // 申请角色工单列表csv
	authModule.POST("/apply_order_list/print", authority.ApplyRoleListPrintServer) // 申请角色工单列表pdf
	authModule.POST("/review_order", authority.ReviewApplyRoleServer)              // 审批申请工单接口
	// 首页数据看版
	dashboardModule := router.Group("/dashboard", CheckLoginStatus)
	dashboardModule.POST("/global_stats", dashboard.GlobalStatsServer)  // 首页全局数据指标
	dashboardModule.POST("/shop_list", dashboard.MapShopInfoListServer) // 首页地图中的店铺信息
	// 识别记录模块
	recordModule := router.Group("/record")
	_ = recordModule
	// 报警记录模块
	alarmModule := router.Group("/alarm", CheckLoginStatus, Authentication)
	alarmModule.POST("/alarm_list", alarm.AlarmListServer)            // 报警记录列表
	alarmModule.POST("/alarm_list/csv", alarm.AlarmListCsvServer)     // 报警记录列表csv
	alarmModule.POST("/alarm_list/print", alarm.AlarmListPrintServer) //报警记录列表pdf
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
	email := ctx.GetHeader("Access-Email")
	token := ctx.GetHeader("Access-Token")
	ctx.Set("email", email)
	if Test {
		ctx.Next()
		return
	}
	key := fmt.Sprintf(util.UserLoginTemplate, email)
	res, err := redis.GetWithRetry(ctx, key)
	if err != nil || res != token {
		util.ErrorJson(ctx, util.UserNotLogin, "用户登录验证失败")
		ctx.AbortWithStatus(http.StatusForbidden)
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
	ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, Access-Token, Access-Email")
	ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	method := ctx.Request.Method
	if method == "OPTIONS" {
		ctx.AbortWithStatus(http.StatusNoContent)
	}
	ctx.Next()
}
