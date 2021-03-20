package tool

import (
	"fsbm/db"
	"fsbm/util"
	"fsbm/util/logs"
	"fsbm/util/mail"
	"github.com/gin-gonic/gin"
)

// 产生邮箱验证码
func GenerateVerificationCode(ctx *gin.Context) {
	var req toolCommonRequest
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
	// 邮箱注册过
	if user == nil {
		logs.CtxInfo(ctx, "email not register")
		util.ErrorJson(ctx, util.UserNotExist, "邮箱未注册")
		return
	}
	code, err := util.SetVerificationCode(ctx, req.Email)
	if err != nil {
		logs.CtxError(ctx, "redis set error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "内部错误")
		return
	}
	err = mail.SendMail(newVerificationMail(req.Email, code))
	if err != nil {
		logs.CtxError(ctx, "send mail error. err: %+v", err)
		util.ErrorJson(ctx, util.EmailSendError, "邮件发送失败")
		return
	}
	util.EndJson(ctx, nil)
}

func newVerificationMail(dest, code string) *mail.DefaultMail {
	return &mail.DefaultMail{
		Dest:    []string{dest},
		Subject: "验证码",
		Text:    []byte(code + "，有效期3分钟"),
	}
}
