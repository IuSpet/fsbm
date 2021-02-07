package userAccount

import (
	"fsbm/db"
	"fsbm/util"
	"fsbm/util/logs"
	"github.com/gin-gonic/gin"
	"regexp"
)

var legalEmailAddr = regexp.MustCompile("")

func UserLoginServer(ctx *gin.Context) {
	var req userCommonRequest
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	db.GetUserByEmail(req.Email)

}

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
	if res.Email != "" {
		logs.CtxInfo(ctx, "email used")
		util.ErrorJson(ctx, util.RepeatedEmailAddr, "邮箱已被注册")
		return
	}
	//TODO 其他参数合规检查
	newUser := &db.UserAccountInfo{
		Name:     req.Name,
		Email:    req.Email,
		Status:   0,
		Password: util.Sha256(util.Salt + req.Password),
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

func LogoutServer(ctx *gin.Context) {

}

// 检查邮箱地址是否合法
func isEmailLegal(email string) bool {
	return legalEmailAddr.MatchString(email)
}
