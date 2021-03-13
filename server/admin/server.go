package admin

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
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
	userList, cnt, err := getSortedUserList(&req)
	if err != nil {
		logs.CtxError(ctx, "get user list error.err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	for idx := range userList {
		rsp.UserInfoList = append(rsp.UserInfoList, userInfo{
			Name:      userList[idx].Name,
			Email:     userList[idx].Email,
			Gender:    db.UserGenderMapping[userList[idx].Gender],
			Age:       userList[idx].Age,
			Phone:     userList[idx].Phone,
			CreatedAt: userList[idx].CreatedAt.Format(util.YMDHMS),
			Status:    db.UserGenderMapping[userList[idx].Status],
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
	userList, _, err := getSortedUserList(&req)
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
	rsp.Status = db.UserStatusMapping[user.Status]
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
		logs.CtxError(ctx, "get user by email error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "内部错误")
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

func getSortedUserList(req *getUserListRequest) ([]db.UserAccountInfo, int64, error) {
	begin, err := time.Parse(util.YMDHMS, req.CreateBegin)
	if err != nil {
		fmt.Printf("1 %+v\n", err)
		begin, err = time.Parse(util.H5FMT, req.CreateBegin)
		fmt.Printf("2 %+v\n", err)
	}
	end, err := time.Parse(util.YMDHMS, req.CreateEnd)
	if err != nil {
		fmt.Printf("3 %+v\n", err)
		end, err = time.Parse(util.H5FMT, req.CreateEnd)
		fmt.Printf("4 %+v\n", err)
	}
	userList, totalCnt, err := getUserList(req.Name, req.Email, req.Phone, req.Gender, req.Age, begin, end, req.Page, req.PageSize)
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
				asc := strings.ToLower(item.Order) == "asc"
				switch x.(type) {
				case string:
					if asc {
						return x.(string) > y.(string)
					}
					return x.(string) < y.(string)
				case int64:
					if asc {
						return x.(int64) > y.(int64)
					}
					return x.(int64) < y.(int64)
				case int8:
					if asc {
						return x.(int8) > y.(int8)
					}
					return x.(int8) < y.(int8)
				case time.Time:
					if asc {
						return x.(time.Time).After(y.(time.Time))
					}
					return x.(time.Time).Before(y.(time.Time))
				}
			}
			return true
		})
	}
	return userList, totalCnt, nil
}
