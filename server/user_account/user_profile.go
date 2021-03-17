package userAccount

import (
	"fmt"
	"fsbm/db"
	"fsbm/util"
	"fsbm/util/logs"
	"fsbm/util/redis"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
)

func SetAvatarServer(ctx *gin.Context) {
	bodyReader := ctx.Request.Body
	data, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		logs.CtxError(ctx, "read body error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "请求内容读取失败")
		return
	}
	if len(data) > 64*64 {
		logs.CtxError(ctx, "body too large. len: %d", len(data))
		util.ErrorJson(ctx, util.ParamError, "内容太大")
		return
	}
	email := ctx.GetHeader("email")
	err = db.SetAvatar(email, data)
	if err != nil {
		logs.CtxError(ctx, "db error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
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
	logs.CtxInfo(ctx, "req: %+v", req)
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
		Gender:   req.Gender,
		Age:      req.Age,
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
	logs.CtxInfo(ctx, "req: %+v", req)
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

// 获取用户头像
func GetAvatarServer(ctx *gin.Context) {

}
