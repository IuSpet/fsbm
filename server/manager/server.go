package manager

import (
	"fsbm/db"
	"fsbm/util"
	"fsbm/util/logs"
	"github.com/gin-gonic/gin"
)

// 获取所有角色列表
func GetRoleListServer(ctx *gin.Context) {
	roleList, err := db.GetRoleList()
	if err != nil {
		logs.CtxError(ctx, "get role list error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	var rsp getRoleListResponse
	for _, role := range roleList {
		rsp.List = append(rsp.List, roleInfo{
			Role: role.Role,
			Id:   role.ID,
		})
	}
	util.EndJson(ctx, rsp)
}

// 获取用户激活角色和过期角色
func GetUserRoleListServer(ctx *gin.Context) {
	var req getUserRoleListRequest
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		logs.CtxError(ctx, "get user info error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	if user == nil {
		logs.CtxWarn(ctx, "user not exist.")
		util.ErrorJson(ctx, util.UserNotExist, "用户不存在")
		return
	}
	userRoleList, err := getUserRoleList(user.ID)
	if err != nil {
		logs.CtxError(ctx, "get user role list error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	var rsp getUserRoleListResponse
	for _, role := range userRoleList {
		switch role.Status {
		case db.AuthUserRoleStatus_Active:
			rsp.ActiveRoles = append(rsp.ActiveRoles, roleInfo{
				Role: role.Role,
				Id:   role.RoleId,
			})
		case db.AuthUserRoleStatus_Expired:
			rsp.ExpiredRoles = append(rsp.ExpiredRoles, roleInfo{
				Role: role.Role,
				Id:   role.RoleId,
			})
		}
	}
	util.EndJson(ctx, rsp)
}

func ApplyRoleServer(ctx *gin.Context) {
	var req applyRoleRequest
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		logs.CtxError(ctx, "get user info error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	if user == nil {
		logs.CtxWarn(ctx, "user not exist.")
		util.ErrorJson(ctx, util.UserNotExist, "用户不存在")
		return
	}
	role, err := db.GetRoleById(req.RoleId)
	if err != nil {
		logs.CtxError(ctx, "get role info error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	row := &db.AuthApplyRole{
		UserId: user.ID,
		Email:  user.Email,
		RoleId: role.ID,
		Role:   role.Role,
		Status: 0,
	}
	err = db.SaveAthApplyRoleRow(row)
	if err != nil {
		logs.CtxError(ctx, "save auth apply role row error. err:%+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	util.EndJson(ctx, nil)
}
