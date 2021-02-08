package userAccount

import (
	"context"
	"fmt"
	"fsbm/db"
	"fsbm/util"
	"fsbm/util/logs"
	"fsbm/util/redis"
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
	"time"
)

const loginExpiration = 30 * time.Minute

var legalEmailAddr = regexp.MustCompile("")

// 密码登陆接口
func UserPasswordLoginServer(ctx *gin.Context) {
	var req userCommonRequest
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	res, err := db.GetUserByEmail(req.Email)
	if err != nil {
		logs.CtxError(ctx, "get user by email error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "内部错误")
		return
	}
	if res == nil {
		logs.CtxInfo(ctx, "user not exist")
		util.ErrorJson(ctx, util.UserNotExist, "该邮箱未注册")
		return
	}
	if res.Status == 1 {
		logs.CtxInfo(ctx, "user has been deleted")
		util.ErrorJson(ctx, util.UserDeleted, "用户已被删除 ")
	}
	if res.Password != encryptPassword(req.Password) {
		logs.CtxInfo(ctx, "wrong password")
		util.ErrorJson(ctx, util.InvalidPassword, "密码错误")
		return
	}
	// 登陆成功
	err = setLoginStatus(ctx, req.Email)
	if err != nil {
		logs.CtxError(ctx, "set login status error")
	}
	util.EndJson(ctx, nil)
}

// 验证登陆接口
func UserVerifyLoginServer(ctx *gin.Context) {
	var req userCommonRequest
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	// 检查用户是否注册
	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		logs.CtxError(ctx, "get user by email error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "内部错误")
		return
	}
	if user == nil {
		logs.CtxInfo(ctx, "user not exist")
		util.ErrorJson(ctx, util.UserNotExist, "该邮箱未注册")
		return
	}
	if user.Status == 1 {
		logs.CtxInfo(ctx, "user has been deleted")
		util.ErrorJson(ctx, util.UserDeleted, "用户已被删除 ")
	}
	// 查询验证码
	key := fmt.Sprintf(util.UserLoginVerificationCodeTemplate, req.Email)
	res, err := redis.GetWithRetry(ctx, key)
	if err != nil {
		logs.CtxWarn(ctx, "redis get error. key: %+v, err: %+v", key, err)
		util.ErrorJson(ctx, util.DbError, "获取验证码失败")
		return
	}
	if res != req.VerifyCode {
		logs.CtxInfo(ctx, "verification error")
		util.ErrorJson(ctx, util.InvalidVerificationCode, "验证码错误")
		return
	}
	err = setLoginStatus(ctx, req.Email)
	if err != nil {
		logs.CtxError(ctx, "set login status error")
	}
	util.EndJson(ctx, nil)
}

// 注册接口
func UserRegisterServer(ctx *gin.Context) {
	var req userCommonRequest
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	if !isEmailLegal(req.Email) {
		logs.CtxInfo(ctx, "illegal email email: %s", req.Email)
		util.ErrorJson(ctx, util.IllegalEmailAddr, "邮箱格式错误")
		return
	}
	res, err := db.GetUserByEmail(req.Email)
	if err != nil {
		logs.CtxError(ctx, "get user by email error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "内部错误")
		return
	}
	// 邮箱注册过
	if res != nil {
		logs.CtxInfo(ctx, "email used")
		util.ErrorJson(ctx, util.RepeatedEmailAddr, "邮箱已被注册")
		return
	}
	//TODO 其他参数合规检查
	newUser := &db.UserAccountInfo{
		Name:     req.Name,
		Email:    req.Email,
		Status:   0,
		Password: encryptPassword(req.Password),
	}
	err = db.SaveUserInfo(newUser)
	if err != nil {
		logs.CtxError(ctx, "save user error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "内部错误")
		return
	}
	logs.CtxInfo(ctx, "new user: %+v", newUser)
	util.EndJson(ctx, nil)
}

// 注销接口
func LogoutServer(ctx *gin.Context) {
	var req userCommonRequest
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	delLoginStatus(ctx, req.Email)
	util.EndJson(ctx, nil)
}

// 修改接口
func ModifyServer(ctx *gin.Context) {
	var req userCommonRequest
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	// 验证登陆
	if !checkLoginStatus(ctx, req.Email) {
		logs.CtxInfo(ctx, "未登陆")
		util.ErrorJson(ctx, util.UserNotLogin, "用户未登陆")
		return
	}
	key := fmt.Sprintf(util.UserLoginVerificationCodeTemplate, req.Email)
	res, err := redis.GetWithRetry(ctx, key)
	if err != nil {
		logs.CtxWarn(ctx, "redis get error. key: %+v, err: %+v", key, err)
		util.ErrorJson(ctx, util.DbError, "获取验证码失败")
		return
	}
	if strings.ToLower(res) != strings.ToLower(req.VerifyCode) {
		logs.CtxInfo(ctx, "verification error, %+v, %+v", res, req.VerifyCode)
		util.ErrorJson(ctx, util.InvalidVerificationCode, "验证码错误")
		return
	}
	existInfo, err := db.GetUserByEmail(req.Email)
	if err != nil {
		logs.CtxError(ctx, "get user by email error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "内部错误")
		return
	}
	modifyInfo := &db.UserAccountInfo{
		ID:       existInfo.ID,
		Name:     req.Name,
		Email:    existInfo.Email,
		Status:   0,
		Password: encryptPassword(req.Password),
	}
	err = db.SaveUserInfo(modifyInfo)
	if err != nil {
		logs.CtxError(ctx, "save user info error, err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "内部错误")
		return
	}
	util.EndJson(ctx, nil)
}

// 删除接口
func DeleteServer(ctx *gin.Context) {
	var req userCommonRequest
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	// 验证登陆
	if !checkLoginStatus(ctx, req.Email) {
		logs.CtxInfo(ctx, "未登陆")
		util.ErrorJson(ctx, util.UserNotLogin, "用户未登陆")
		return
	}
	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		logs.CtxError(ctx, "get user by email error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "内部错误")
		return
	}
	user.Status = 1
	err = db.SaveUserInfo(user)
	if err != nil {
		logs.CtxError(ctx, "save user info error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "内部错误")
		return
	}
	util.EndJson(ctx, nil)
}

// 检查邮箱地址是否合法
func isEmailLegal(email string) bool {
	return legalEmailAddr.MatchString(email)
}

// 密码加盐后sha256加密
func encryptPassword(password string) string {
	return util.Sha256(util.Salt + password)
}

// 写入登陆状态
func setLoginStatus(ctx context.Context, email string) error {
	key := fmt.Sprintf(util.UserLoginTemplate, email)
	err := redis.SetWithRetry(ctx, key, "ok", loginExpiration)
	return err
}

// 删除登陆状态
func delLoginStatus(ctx context.Context, email string) {
	key := fmt.Sprintf(util.UserLoginTemplate, email)
	redis.Del(ctx, key)
}

func checkLoginStatus(ctx context.Context, email string) bool {
	key := fmt.Sprintf(util.UserLoginTemplate, email)
	res, err := redis.GetWithRetry(ctx, key)
	if err != nil || res == "" {
		return false
	}
	return true
}
