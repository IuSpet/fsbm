package authority

import (
	"errors"
	"fmt"
	"fsbm/db"
	"fsbm/util"
	"fsbm/util/logs"
	"github.com/gin-gonic/gin"
	"reflect"
	"sort"
	"strings"
	"time"
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

// 申请角色请求提交
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

// 用户请求工单列表
func ApplyRoleListServer(ctx *gin.Context) {
	req := newApplyRoleListRequest()
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	applyOrderList, totalCnt, err := getSortedApplyOrderList(req, false)
	if err != nil {
		logs.CtxError(ctx, "get apply order list error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	var rsp applyRoleListResponse
	for _, row := range applyOrderList {
		rsp.List = append(rsp.List, applyRoleOrder{
			Id:           row.Id,
			User:         row.User,
			Role:         row.Role,
			Reason:       row.Reason,
			Status:       db.AuthApplyRoleStatusMapping[row.Status],
			Reviewer:     row.Reviewer,
			ReviewReason: row.ReviewReason,
			ReviewAt:     time.Unix(row.ReviewAt, 0).Format(util.YMDHMS),
		})
	}
	rsp.TotalCnt = totalCnt
	util.EndJson(ctx, rsp)
}

func newApplyRoleListRequest() *applyRoleListRequest {
	return &applyRoleListRequest{
		User:      "",
		Role:      "",
		Reviewer:  "",
		Status:    []int8{1},
		BeginDate: time.Unix(0, 0).Format(util.YMDHMS),
		EndDate:   time.Now().Format(util.YMDHMS),
		ListReqField: util.ListReqField{
			Page:       1,
			PageSize:   10,
			SortFields: nil,
		},
	}
}

func getSortedApplyOrderList(req *applyRoleListRequest, all bool) ([]applyRoleRow, int64, error) {
	begin, err1 := time.Parse(util.YMDHMS, req.BeginDate)
	end, err2 := time.Parse(util.YMDHMS, req.EndDate)
	if err1 != nil || err2 != nil {
		return nil, 0, errors.New(fmt.Sprintf("err1: %+v, err2: %+v", err1, err2))
	}
	applyOrderList, err := getApplyRoleOrderList(req.User, req.Role, req.Reviewer, req.Status, begin.Unix(), end.Unix())
	if err != nil {
		return nil, 0, err
	}
	if len(req.SortFields) > 0 {
		sort.SliceStable(applyOrderList, func(i, j int) bool {
			a, b := reflect.ValueOf(applyOrderList[i]), reflect.ValueOf(applyOrderList[j])
			for _, item := range req.SortFields {
				x, y := a.FieldByName(item.Field).Interface(), b.FieldByName(item.Field).Interface()
				if reflect.DeepEqual(x, y) {
					continue
				}
				desc := strings.ToLower(item.Order) == "desc"
				less := util.CmpInterface(x, y)
				if desc {
					return !less
				}
				return less
			}
			return true
		})
	}
	totalCnt := int64(len(applyOrderList))
	if all {
		return applyOrderList, totalCnt, nil
	}
	offset := req.PageSize * (req.Page - 1)
	if offset+req.PageSize >= totalCnt {
		return applyOrderList[offset:], totalCnt, nil
	}
	return applyOrderList[offset : offset+req.PageSize], totalCnt, nil
}
