package admin

import (
	"encoding/csv"
	"encoding/json"
	"fsbm/db"
	"fsbm/util"
	"fsbm/util/logs"
	"github.com/gin-gonic/gin"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

var userInfoCsvColumns = []struct{ Field, Key string }{
	{"用户名", "name"},
	{"邮箱", "email"},
	{"手机号", "phone"},
	{"性别", "gender"},
	{"年龄", "age"},
	{"注册时间", "created_at"},
	{"状态", "status"},
}

// 获取所有用户列表接口
func UserListServer(ctx *gin.Context) {
	req := newGetUserListRequest()
	var rsp getUserListResponse
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	userList, cnt, err := getSortedUserList(&req, false)
	if err != nil {
		logs.CtxError(ctx, "get user list error.err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	for idx := range userList {
		rsp.UserInfoList = append(rsp.UserInfoList, userInfo{
			Id:        userList[idx].ID,
			Name:      userList[idx].Name,
			Email:     userList[idx].Email,
			Gender:    db.UserGenderMapping[userList[idx].Gender],
			Age:       userList[idx].Age,
			Phone:     userList[idx].Phone,
			CreatedAt: userList[idx].CreatedAt.Format(util.YMDHMS),
			Status:    db.UserStatusMapping[userList[idx].Status],
		})
	}
	rsp.TotalCount = cnt
	util.EndJson(ctx, rsp)
}

// 查询用户信息导出csv接口
func UserListCsvServer(ctx *gin.Context) {
	req := newGetUserListRequest()
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	req.Page = -1
	logs.CtxInfo(ctx, "req: %+v", req)
	userList, _, err := getSortedUserList(&req, true)
	if err != nil {
		logs.CtxError(ctx, "get user list error.err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	csvRows := make([]userInfoCsvRow, 0, len(userList))
	for idx := range userList {
		csvRows = append(csvRows, userInfoCsvRow{
			Email:     userList[idx].Email,
			Name:      userList[idx].Name,
			Phone:     userList[idx].Phone,
			Status:    db.UserStatusMapping[userList[idx].Status],
			Age:       strconv.FormatInt(int64(userList[idx].Age), 10),
			CreatedAt: userList[idx].CreatedAt.Format(util.YMDHMS),
			Gender:    db.UserGenderMapping[userList[idx].Gender],
		})
	}
	fileName := "用户列表导出.csv"
	file, _ := os.Create(fileName)
	defer file.Close()
	w := csv.NewWriter(file)
	_, _ = file.WriteString("\xEF\xBB\xBF")
	var title []string
	for idx := range userInfoCsvColumns {
		title = append(title, userInfoCsvColumns[idx].Field)
	}
	_ = w.Write(title)
	for idx := range csvRows {
		var row []string
		var m = make(map[string]string)
		s, _ := json.Marshal(csvRows[idx])
		_ = json.Unmarshal(s, &m)
		for _, k := range userInfoCsvColumns {
			row = append(row, m[k.Key])
		}
		_ = w.Write(row)
	}
	w.Flush()
	util.SetFileTransportHeader(ctx, fileName)
	_ = os.Remove(fileName)
}

// 查询用户信息，制作打印界面
func UserListPrintServer(ctx *gin.Context) {
	req := newGetUserListRequest()
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	req.Page = -1
	logs.CtxInfo(ctx, "req: %+v", req)
	userList, _, err := getSortedUserList(&req, true)
	if err != nil {
		logs.CtxError(ctx, "get user list error.err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	var rsp getUserListResponse
	for _, item := range userList {
		rsp.UserInfoList = append(rsp.UserInfoList, userInfo{
			Id:        item.ID,
			Name:      item.Name,
			Email:     item.Email,
			Gender:    db.UserGenderMapping[item.Gender],
			Age:       item.Age,
			Phone:     item.Phone,
			CreatedAt: item.CreatedAt.Format(util.YMDHMS),
			Status:    db.UserStatusMapping[item.Status],
		})
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
	user, err := db.GetUserById(req.UserId)
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
	rsp.Email = user.Email
	rsp.Status = db.UserStatusMapping[user.Status]
	rsp.Name = user.Name
	rsp.Age = int64(user.Age)
	rsp.Gender = db.UserGenderMapping[user.Gender]
	rsp.Phone = user.Phone
	rsp.CreatedAt = user.CreatedAt.Format(util.YMDHMS)
	rsp.Id = user.ID
	util.EndJson(ctx, rsp)
}

// 修改用户信息接口
func ModifyUserDetailServer(ctx *gin.Context) {
	var req modifyUserDetailRequest
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	// 修改用户基本信息
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
	user.Name = req.Name
	user.Status = req.Status
	err = db.SaveUserInfo(user)
	if err != nil {
		logs.CtxError(ctx, "save user info error, err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "内部错误")
		return
	}
	// 增加用户角色
	msg := ""
	if len(req.AddRoles) > 0 {
		err = db.SaveAuthUserRoleRows(generateAuthUserRoleRows(user.ID, req.AddRoles))
		if err != nil {
			logs.CtxWarn(ctx, "user add role fail. err: %+v", err)
			msg += "用户角色增加失败"
		}
	}
	// 删除用户角色
	if len(req.DeleteRoles) > 0 {
		err = db.RemoveUserRole(user.ID, req.DeleteRoles)
		if err != nil {
			logs.CtxWarn(ctx, "user delete role fail. err: %+v", err)
			msg += "用户角色删除失败"
		}
	}
	if msg != "" {
		util.ErrorJson(ctx, util.DbError, msg)
		return
	}
	util.EndJson(ctx, nil)
}

// 用户注册统计接口
func GetUserRegisterInfoServer(ctx *gin.Context) {
	req := newGetUserListRequest()
	var rsp userRegisterInfoResponse
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	begin, _ := time.Parse(util.YMDHMS, req.CreateBegin)
	end, _ := time.Parse(util.YMDHMS, req.CreateEnd)
	userList, err := getUserList(req.Name, req.Email, req.Phone, req.Gender, req.Age, begin, end)
	if err != nil {
		logs.CtxError(ctx, "get user list error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	registerStats := make(map[string]int64)
	for _, row := range userList {
		registerStats[row.CreatedAt.Format(util.YMD)] += 1
	}
	end = end.AddDate(0, 0, 1)
	for begin.Before(end) {
		cur := begin.Format(util.YMD)
		rsp.Series = append(rsp.Series, registerInfo{
			Date: cur,
			Cnt:  registerStats[cur],
		})
		begin = begin.AddDate(0, 0, 1)
	}
	util.EndJson(ctx, rsp)
}

// 为用户增加角色
func AddUserRoleServer(ctx *gin.Context) {
	req := newAddUserRoleRequest()
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	if req.UserId == -1 || req.RoleId == -1 {
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	userRole, err := db.GetUserRoleRow(req.UserId, req.RoleId)
	if err != nil {
		logs.CtxError(ctx, "get user role error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	if userRole == nil {
		userRole = &db.AuthUserRole{
			UserID:    req.UserId,
			RoleID:    req.RoleId,
			StartTime: time.Now(),
			EndTime:   time.Now().Add(time.Duration(req.Expire) * time.Second),
			Status:    db.AuthUserRoleStatus_Active,
		}
	}
	err = db.SaveAuthUserRoleRow(userRole)
	if err != nil {
		logs.CtxError(ctx, "save user role error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	util.EndJson(ctx, nil)
}

// 删除用户角色接口
func DeleteUserRoleServer(ctx *gin.Context) {
	req := newDeleteUserRoleRequest()
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	if req.UserId == -1 || req.RoleId == -1 {
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	userRole, err := db.GetUserRoleRow(req.UserId, req.RoleId)
	if err != nil {
		logs.CtxError(ctx, "get user role error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	if userRole == nil {
		util.EndJson(ctx, nil)
		return
	}
	userRole.Status = db.AuthUserRoleStatus_Expired
	err = db.SaveAuthUserRoleRow(userRole)
	if err != nil {
		logs.CtxError(ctx, "save user role error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	util.EndJson(ctx, nil)
}

// 用户操作记录接口
func GetUserOperationListServer(ctx *gin.Context) {
	req := &getUserOperationListRequest{UserId: -1}
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	rows, err := db.GetUserOperationRows(req.UserId)
	if err != nil {
		logs.CtxError(ctx, "get user operation list error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	rsp := getUserOperationListResponse{}
	for _, row := range rows {
		rsp.List = append(rsp.List, userOperation{
			Operation:  row.Operation,
			OperatedAt: time.Unix(row.OperatedAt, 0).Format(util.YMDHMS),
		})
	}
	util.EndJson(ctx, rsp)
}

func ModifyUserInfoServer(ctx *gin.Context) {
	req := &modifyUserInfoRequest{}
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	user, err := db.GetUserById(req.Id)
	if err != nil || user == nil {
		logs.CtxError(ctx, "get user error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	user.Name = req.Name
	user.Age = req.Age
	user.Status = req.Status
	user.Gender = req.Gender
	if req.Password != "" {
		user.Password = encryptPassword(req.Password)
	}
	err = db.SaveUserInfo(user)
	if err != nil {
		logs.CtxError(ctx, "save user error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	util.EndJson(ctx, nil)
}

func generateAuthUserRoleRows(userID int64, roleIDList []int64) []db.AuthUserRole {
	var userRoleList []db.AuthUserRole
	for _, roleID := range roleIDList {
		userRoleList = append(userRoleList, db.AuthUserRole{
			UserID: userID,
			RoleID: roleID,
			Status: 1,
		})
	}
	return userRoleList
}

// 默认值
func newGetUserListRequest() getUserListRequest {
	return getUserListRequest{
		Gender:      -1,
		Age:         -1,
		CreateBegin: time.Unix(0, 0).Format(util.YMDHMS),
		CreateEnd:   time.Now().Format(util.YMDHMS),
		Page:        1,
		PageSize:    20,
	}
}

func newRegisterStatsRequest() *registerStatsRequest {
	now := time.Now().Add(24 * time.Hour)
	return &registerStatsRequest{
		CreateBegin: now.AddDate(0, 0, -14).Format(util.YMD),
		CreateEnd:   now.Format(util.YMD),
	}
}

func getSortedUserList(req *getUserListRequest, all bool) ([]db.UserAccountInfo, int64, error) {
	begin, _ := time.Parse(util.YMDHMS, req.CreateBegin)
	end, _ := time.Parse(util.YMDHMS, req.CreateEnd)
	userList, err := getUserList(req.Name, req.Email, req.Phone, req.Gender, req.Age, begin, end)
	if err != nil {
		return nil, 0, err
	}
	if len(req.SortFields) > 0 {
		sort.SliceStable(userList, func(i, j int) bool {
			a, b := reflect.ValueOf(userList[i]), reflect.ValueOf(userList[j])
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
	totalCnt := int64(len(userList))
	if all {
		return userList, totalCnt, nil
	}
	offset := (req.Page - 1) * req.PageSize
	if totalCnt > offset+req.PageSize {
		return userList[offset : offset+req.PageSize], totalCnt, nil
	} else if totalCnt > offset {
		return userList[offset:], totalCnt, nil
	} else {
		return nil, totalCnt, nil
	}
}

func newAddUserRoleRequest() *addUserRoleRequest {
	return &addUserRoleRequest{
		UserId: -1,
		RoleId: -1,
		Expire: 30 * 24 * 3600,
	}
}

func newDeleteUserRoleRequest() *deleteUserRoleRequest {
	return &deleteUserRoleRequest{
		UserId: -1,
		RoleId: -1,
	}
}

// 密码加盐后sha256加密
func encryptPassword(password string) string {
	return util.Sha256(util.Salt + password)
}
