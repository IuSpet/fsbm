package tool

import (
	"bytes"
	"fmt"
	"fsbm/conf"
	"fsbm/db"
	"fsbm/util"
	"fsbm/util/logs"
	"fsbm/util/mail"
	"github.com/gin-gonic/gin"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"time"
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
	// 没有该邮箱信息
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

// 保存图片
func SaveImgServer(ctx *gin.Context) {
	bodyReader := ctx.Request.Body
	data, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		logs.CtxError(ctx, "read body error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "请求内容读取失败")
		return
	}
	email := ctx.GetString("email")
	path, err := saveImgLocal(data, email)
	if err != nil {
		logs.CtxError(ctx, "save img error. err: %+v", err)
		util.ErrorJson(ctx, util.SaveImgError, "图片保存失败")
		return
	}
	util.EndJson(ctx, saveImgResponse{Path: path})
}

// 生成发送验证码的邮件
func newVerificationMail(dest, code string) *mail.DefaultMail {
	return &mail.DefaultMail{
		Dest:    []string{dest},
		Subject: "验证码",
		Text:    []byte(code + "，有效期3分钟"),
	}
}

// 保存图片到本地
func saveImgLocal(blob []byte, email string) (string, error) {
	now := time.Now()
	name := fmt.Sprintf("%s%d.png", email, now.Unix())
	var path string
	env := conf.GetEnv()
	switch env {
	case conf.TEST:
		path = "~/fsbm_test/img/"
	case conf.PRODUCT:
		path = "/fsbm/img/"
	}
	file, err := os.Create(path + name)
	if err != nil {
		return "", err
	}
	img, _, err := image.Decode(bytes.NewReader(blob))
	if err != nil {
		return "", err
	}
	err = png.Encode(file, img)
	if err != nil {
		return "", err
	}
	return path + name, nil
}
