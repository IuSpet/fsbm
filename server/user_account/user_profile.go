package userAccount

import (
	"bytes"
	"fmt"
	"fsbm/db"
	"fsbm/util"
	"fsbm/util/logs"
	"github.com/gin-gonic/gin"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func SetAvatarServer(ctx *gin.Context) {
	//logs.CtxInfo(ctx, "header: %+v", ctx.Request.Header)
	bodyReader := ctx.Request.Body
	data, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		logs.CtxError(ctx, "read body error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "请求内容读取失败")
		return
	}
	if len(data) > 64*64*4 {
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
	code, err := util.GetVerificationCode(ctx, req.Email)
	if err != nil {
		logs.CtxWarn(ctx, "redis get error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "获取验证码失败")
		return
	}
	if strings.ToLower(code) != strings.ToLower(req.VerifyCode) {
		logs.CtxInfo(ctx, "verification error, %+v, %+v", code, req.VerifyCode)
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
	var req deleteUserRequest
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	code, err := util.GetVerificationCode(ctx, req.Email)
	if err != nil {
		logs.CtxWarn(ctx, "redis get error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "获取验证码失败")
		return
	}
	if strings.ToLower(code) != strings.ToLower(req.VerifyCode) {
		logs.CtxInfo(ctx, "verification error, %+v, %+v", code, req.VerifyCode)
		util.ErrorJson(ctx, util.InvalidVerificationCode, "验证码错误")
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

// 获取用户信息
func GetUserProfile(ctx *gin.Context) {
	var req getUserProfileRequest
	logs.CtxInfo(ctx, "header: %+v", ctx.Request.Header)
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		logs.CtxError(ctx, "get user by email error. err: %+v")
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	if user == nil {
		logs.CtxInfo(ctx, "user not exist.")
		util.ErrorJson(ctx, util.UserNotExist, "用户不存在")
		return
	}
	rsp := getUserProfileResponse{
		Name:      user.Name,
		Phone:     user.Phone,
		Gender:    db.UserGenderMapping[user.Gender],
		Age:       user.Age,
		CreatedAt: user.CreatedAt.Format(util.YMD),
	}
	util.EndJson(ctx, rsp)
}

// 获取用户头像
func GetAvatarServer(ctx *gin.Context) {
	var req getUserProfileRequest
	logs.CtxInfo(ctx, "body: %+v", ctx.Request.Body)
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		logs.CtxError(ctx, "get user by email error. err: %+v")
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	if user == nil {
		logs.CtxInfo(ctx, "user not exist.")
		util.ErrorJson(ctx, util.UserNotExist, "用户不存在")
		return
	}
	avatar := util.NewAvatarHandler(user.Avatar)
	if avatar.ContentLength == 0 {
		util.ErrorJson(ctx, util.AvatarNotExist, "未设置头像")
		return
	}
	//avatar.SetHeaders("Content-Disposition", `attachment;filename="avatar.png"`)
	ctx.DataFromReader(http.StatusOK, avatar.ContentLength, avatar.ContentType, avatar.Avatar, avatar.ExtraHeaders)
}

// 获取用户角色
func GetUserRolesServer(ctx *gin.Context) {
	var req getUserProfileRequest
	logs.CtxInfo(ctx, "body: %+v", ctx.Request.Body)
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		logs.CtxError(ctx, "get user by email error. err: %+v")
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	if user == nil {
		logs.CtxInfo(ctx, "user not exist.")
		util.ErrorJson(ctx, util.UserNotExist, "用户不存在")
		return
	}
	roles, err := db.GetRoleByUserId(user.ID)
	var rsp getUserRolesResponse
	for _, role := range roles {
		rsp.Roles = append(rsp.Roles, role.Role)
	}
	util.EndJson(ctx, rsp)
}

func saveImg(b []byte) {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		fmt.Println(err)
	}
	out, _ := os.Create("./img.png")
	err = png.Encode(out, img)
	if err != nil {
		fmt.Println(err)
	}
}
