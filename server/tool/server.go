package tool

import (
	"fmt"
	"fsbm/util"
	"fsbm/util/logs"
	"fsbm/util/mail"
	"fsbm/util/redis"
	"github.com/gin-gonic/gin"
	"time"
)

const verificationExpiration = time.Minute

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
	code := util.GenerateRandCode(6)
	key := fmt.Sprintf(util.UserLoginVerificationCodeTemplate, req.Email)
	err = redis.SetWithRetry(ctx, key, code, verificationExpiration)
	if err != nil {
		logs.CtxError(ctx, "redis set error. err: %+v", err)
		return
	}
	err = mail.SendMail(newVerificationMail(req.Email, code))
	if err != nil {
		logs.CtxError(ctx, "send mail error. err: %+v", err)
		return
	}
}

func newVerificationMail(dest, code string) *mail.DefaultMail {
	return &mail.DefaultMail{
		Dest:    []string{dest},
		Subject: "登陆验证码",
		Text:    []byte(code),
	}
}
