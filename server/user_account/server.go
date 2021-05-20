package userAccount

import (
	"context"
	"fmt"
	"fsbm/db"
	"fsbm/task"
	"fsbm/util"
	"fsbm/util/logs"
	"fsbm/util/redis"
	"github.com/gin-gonic/gin"
	"regexp"
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
	if req.Email == "" || req.Password == "" {
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
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
		return
	}
	if res.Password != encryptPassword(req.Password) {
		logs.CtxInfo(ctx, "wrong password")
		util.ErrorJson(ctx, util.InvalidPassword, "密码错误")
		return
	}
	// 登陆成功
	token := util.Md5(time.Now().Format(util.YMDHMS) + req.Email)
	err = setLoginStatus(ctx, req.Email, token)
	if err != nil {
		logs.CtxError(ctx, "set login status error")
	}
	util.EndJson(ctx, loginResponse{
		Email: req.Email,
		Token: token,
	})
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
	key := fmt.Sprintf(util.UserVerificationCodeTemplate, req.Email)
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
	token := util.Md5(time.Now().Format(util.YMDHMS) + req.Email)
	err = setLoginStatus(ctx, req.Email, token)
	if err != nil {
		logs.CtxError(ctx, "set login status error")
	}
	util.EndJson(ctx, loginResponse{
		Email: req.Email,
		Token: token,
	})
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
		Phone:    req.Phone,
		Age:      req.Age,
		Gender:   req.Gender,
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
	token := util.Md5(time.Now().Format(util.YMDHMS) + req.Email)
	err = setLoginStatus(ctx, req.Email, token)
	if err != nil {
		logs.CtxError(ctx, "set login status error")
	}
	_ = task.SetUserOperationById(newUser.ID, util.UserOPeration_Register)
	util.EndJson(ctx, loginResponse{
		Email: req.Email,
		Token: token,
	})
}

// 注销接口
func LogoutServer(ctx *gin.Context) {
	email := ctx.GetString("email")
	delLoginStatus(ctx, email)
	util.EndJson(ctx, nil)
}

// 申请角色
func ApplyRoleServer(ctx *gin.Context) {
	var req applyRoleRequest
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
	err = db.SaveUserApplyRoleRows(GenerateApplyRoleRows(user.ID, req.RoleIDList))
	if err != nil {
		logs.CtxError(ctx, "save user apply role rows error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "内部错误")
		return
	}
	util.EndJson(ctx, nil)
}

func GenerateApplyRoleRows(UserID int64, RoleIDList []int64) []db.UserApplyRole {
	var rows []db.UserApplyRole
	for _, roleID := range RoleIDList {
		rows = append(rows, db.UserApplyRole{
			UserID: UserID,
			RoleID: roleID,
			Status: 1,
		})
	}
	return rows
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
func setLoginStatus(ctx context.Context, email, token string, ) error {
	key := fmt.Sprintf(util.UserLoginTemplate, email)
	logs.CtxInfo(ctx, key)
	err := redis.SetWithRetry(ctx, key, token, loginExpiration)
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
