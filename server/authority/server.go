package authority

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"fsbm/db"
	"fsbm/task"
	"fsbm/util"
	"fsbm/util/logs"
	"github.com/gin-gonic/gin"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"
)

var applyOrderColumns = []struct{ Field, Key string }{
	{"申请用户", "user"},
	{"申请角色", "role"},
	{"申请理由", "reason"},
	{"工单状态", "status"},
	{"审批人", "reviewer"},
	{"审批理由", "review_reason"},
	{"审批时间", "review_at"},
	{"工单提交时间", "created_at"},
}

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
	req := newGetUserRoleListRequest()
	err := ctx.Bind(req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	var user *db.UserAccountInfo
	if req.UserId != -1 {
		user, err = db.GetUserById(req.UserId)
	} else if req.Email != "" {
		user, err = db.GetUserByEmail(req.Email)
	} else {
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
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
	logs.CtxDebug(ctx, "rsp: %+v", rsp)
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
	logs.CtxInfo(ctx, "req: %+v", req)
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
		UserId:       user.ID,
		Email:        user.Email,
		RoleId:       role.ID,
		Role:         role.Role,
		Expiration:   req.Expiration,
		Reason:       req.Reason,
		ReviewUserId: -1,
		Status:       db.AuthApplyRoleStatus_Unreviewd,
	}
	err = db.SaveAuthApplyRoleRow(row)
	if err != nil {
		logs.CtxError(ctx, "save auth apply role row error. err:%+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	_ = task.SetUserOperationByEmail(ctx.GetString("email"), fmt.Sprintf(util.UserOperation_ApplyRole, role.ID, role.Role))
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
	logs.CtxInfo(ctx, "req: %+v", req)
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
			ReviewAt:     getReviewTime(row.ReviewAt),
			CreatedAt:    row.CreatedAt.Format(util.YMDHMS),
		})
	}
	rsp.TotalCnt = totalCnt
	util.EndJson(ctx, rsp)
}

// 用户申请工单列表csv
func ApplyRoleListCsvServer(ctx *gin.Context) {
	req := newApplyRoleListRequest()
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	applyOrderList, _, err := getSortedApplyOrderList(req, true)
	if err != nil {
		logs.CtxError(ctx, "get apply order list error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	csvRows := make([]applyRoleCsvRow, 0, len(applyOrderList))
	for _, row := range applyOrderList {
		csvRows = append(csvRows, applyRoleCsvRow{
			User:         row.User,
			Role:         row.Role,
			Reason:       row.Reason,
			Status:       db.AuthApplyRoleStatusMapping[row.Status],
			Reviewer:     row.Reviewer,
			ReviewReason: row.ReviewReason,
			ReviewAt:     getReviewTime(row.ReviewAt),
			CreatedAt:    row.CreatedAt.Format(util.YMDHMS),
		})
	}
	fileName := "权限申请工单列表导出.csv"
	file, _ := os.Create(fileName)
	defer file.Close()
	w := csv.NewWriter(file)
	_, _ = file.WriteString("\xEF\xBB\xBF")
	title := make([]string, len(applyOrderColumns))
	for _, item := range applyOrderColumns {
		title = append(title, item.Field)
	}
	_ = w.Write(title)
	for idx := range csvRows {
		var row []string
		var m = make(map[string]string)
		s, _ := json.Marshal(csvRows[idx])
		_ = json.Unmarshal(s, &m)
		for _, k := range applyOrderColumns {
			row = append(row, m[k.Key])
		}
		_ = w.Write(row)
	}
	w.Flush()
	util.SetFileTransportHeader(ctx, fileName)
	_ = os.Remove(fileName)
}

func ApplyRoleListPrintServer(ctx *gin.Context) {
	req := newApplyRoleListRequest()
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	applyOrderList, totalCnt, err := getSortedApplyOrderList(req, true)
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
			ReviewAt:     getReviewTime(row.ReviewAt),
			CreatedAt:    row.CreatedAt.Format(util.YMDHMS),
		})
	}
	rsp.TotalCnt = totalCnt
	util.EndJson(ctx, rsp)
}

// 审批用户申请
func ReviewApplyRoleServer(ctx *gin.Context) {
	req := &reviewApplyRoleRequest{}
	now := time.Now()
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	if !isValidReview(req) {
		logs.CtxWarn(ctx, "invalid req")
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	order, err := db.GetAuthApplyRoleById(req.Id)
	if err != nil {
		logs.CtxError(ctx, "get auth apply role order error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	if order == nil {
		logs.CtxWarn(ctx, "order not found.")
		util.ErrorJson(ctx, util.ApplyRoleOrderNotExist, "工单不存在")
		return
	}
	// 判断工单是否已被审批
	if order.Status != db.AuthApplyRoleStatus_Unreviewd {
		logs.CtxInfo(ctx, "order has reviewed.")
		util.ErrorJson(ctx, util.ApplyROleOrderHasReviewed, "工单已被审核")
		return
	}
	// 查询神户用户id
	reviewUser, err := db.GetUserByEmail(req.Reviewer)
	if err != nil {
		logs.CtxError(ctx, "get user info error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	if reviewUser == nil {
		logs.CtxWarn(ctx, "user not exist.")
		util.ErrorJson(ctx, util.UserNotExist, "用户不存在")
		return
	}
	// 修改工单状态
	order.ReviewUserId = reviewUser.ID
	order.ReviewReason = req.Reason
	order.ReviewAt = now.Unix()
	switch req.Review {
	case 0:
		order.Status = db.AuthApplyRoleStatus_Approve
	case 1:
		order.Status = db.AuthApplyRoleStatus_Deny
	}
	err = db.SaveAuthApplyRoleRow(order)
	if err != nil {
		logs.CtxError(ctx, "save user apply role order error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	if order.Status == db.AuthApplyRoleStatus_Deny {
		util.EndJson(ctx, nil)
		return
	}
	// 开通权限
	userRole, err := db.GetUserRoleRow(order.UserId, order.RoleId)
	if err != nil {
		logs.CtxError(ctx, "get user role row error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	// 没有记录新建记录
	if userRole == nil {
		userRole = &db.AuthUserRole{
			UserID:    order.UserId,
			RoleID:    order.RoleId,
			StartTime: time.Time{},
			EndTime:   time.Time{},
			Status:    0,
		}
	}
	userRole.Status = db.AuthUserRoleStatus_Active
	userRole.StartTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	userRole.EndTime = userRole.StartTime.Add(time.Duration(order.Expiration) * time.Second)
	err = db.SaveAuthUserRoleRow(userRole)
	if err != nil {
		logs.CtxError(ctx, "save user role row error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	// 通过所有申请工单
	orderList, err := db.GetAuthApplyRoleByUserRole(order.UserId, order.RoleId)
	if err != nil {
		logs.CtxError(ctx, "get auth apply role rows error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	for i := range orderList {
		orderList[i].ReviewReason = "自动通过"
		orderList[i].ReviewAt = now.Unix()
		orderList[i].ReviewUserId = -1
		orderList[i].Status = db.AuthApplyRoleStatus_Approve
	}
	err = db.SaveAuthApplyRoleRows(orderList)
	if err != nil {
		logs.CtxError(ctx, "save auth apply role rows error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	util.EndJson(ctx, nil)
	return
}

func newApplyRoleListRequest() *applyRoleListRequest {
	return &applyRoleListRequest{
		User:            "",
		Role:            nil,
		Reviewer:        "",
		Status:          []int8{db.AuthApplyRoleStatus_Unreviewd},
		ApplyBeginTime:  time.Unix(0, 0).Format(util.YMDHMS),
		ApplyEndTime:    time.Now().Format(util.YMDHMS),
		ReviewBeginTime: "",
		ReviewEndTime:   "",
		ListReqField: util.ListReqField{
			Page:       1,
			PageSize:   10,
			SortFields: nil,
		},
	}
}

func getSortedApplyOrderList(req *applyRoleListRequest, all bool) ([]applyRoleRow, int64, error) {
	applyBegin, err1 := time.Parse(util.YMDHMS, req.ApplyBeginTime)
	applyEnd, err2 := time.Parse(util.YMDHMS, req.ApplyEndTime)
	if err1 != nil || err2 != nil {
		return nil, 0, errors.New("time param error")
	}
	var reviewBeginAt, reviewEndAt int64
	reviewBeginAt = -1
	reviewEndAt = -1
	if req.ReviewBeginTime != "" && req.ReviewEndTime != "" {
		reviewBegin, err3 := time.Parse(util.YMDHMS, req.ReviewBeginTime)
		reviewEnd, err4 := time.Parse(util.YMDHMS, req.ReviewEndTime)
		if err3 != nil || err4 != nil {
			return nil, 0, errors.New("time param error")
		}
		reviewBeginAt = reviewBegin.Unix()
		reviewEndAt = reviewEnd.Unix()
	}
	applyOrderList, err := getApplyRoleOrderList(req.User, req.Reviewer, req.Role, req.Status, applyBegin, applyEnd, reviewBeginAt, reviewEndAt)
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

func isValidReview(req *reviewApplyRoleRequest) bool {
	if req.Id < 0 {
		return false
	}
	if req.Review != 0 && req.Review != 1 {
		return false
	}
	return true
}

func newGetUserRoleListRequest() *getUserRoleListRequest {
	return &getUserRoleListRequest{
		UserId: -1,
		Email:  "",
	}
}

func getReviewTime(ts int64) string {
	if ts == 0 {
		return ""
	}
	return time.Unix(ts, 0).Format(util.YMDHMS)
}
