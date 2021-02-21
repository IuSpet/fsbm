package admin

import (
	"fsbm/db"
	"fsbm/util"
	"fsbm/util/logs"
	"github.com/gin-gonic/gin"
)

var userStatusMapping = map[int8]string{
	0: "正常",
	1: "已删除",
}

// 获取所有用户列表接口
func UserListServer(ctx *gin.Context) {
	var req getUserListRequest
	var rsp getUserListResponse
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	userList, err := db.GetAllUser()
	if err != nil {
		logs.CtxError(ctx, "get all user error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "内部错误")
		return
	}
	offset := (req.Page - 1) * req.PageSize
	rsp.TotalCount = int64(len(userList))
	if offset < rsp.TotalCount {
		userList = userList[offset:util.MinInt64(rsp.TotalCount, offset+req.PageSize)]
		for _, user := range userList {
			rsp.UserInfoList = append(rsp.UserInfoList, userInfo{
				Name:   user.Name,
				Email:  user.Email,
				Status: userStatusMapping[user.Status],
			})
		}
	}
	util.EndJson(ctx, rsp)
}

// 获取用户详细信息接口
func UserDetailServer(ctx *gin.Context) {
	var req getUserDetailRequest
	var rsp getUserDetailResponse
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
	rsp.Email = user.Email
	rsp.Status = userStatusMapping[user.Status]
	rsp.Name = user.Name
	roleList, err := db.GetRoleById(user.ID)
	if err != nil {
		logs.CtxError(ctx, "get role by id error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "内部错误")
		return
	}
	for _, role := range roleList {
		rsp.Roles = append(rsp.Roles, userDetailRole{
			Type: role.Type,
			Name: role.Role,
		})
	}
	util.EndJson(ctx, rsp)
}
