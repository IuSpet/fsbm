package shop

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

var shopInfoColumns = []struct{ Field, Key string }{
	{"店铺名称", "name"},
	{"负责人", "admin_name"},
	{"负责人电话", "admin_phone"},
	{"负责人邮箱", "admin_email"},
	{"店铺地址", "addr"},
	{"注册时间", "created_at"},
	{"状态", "status"},
}

// 获取店铺列表
func GetShopListServer(ctx *gin.Context) {
	req := newGetShopListRequest()
	err := ctx.Bind(req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	shopInfoList, totalCnt, err := getSortedShopListData(req, false)
	if err != nil {
		logs.CtxError(ctx, "err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	rsp := getShopListResponse{TotalCnt: totalCnt}
	for i := range shopInfoList {
		rsp.List = append(rsp.List, shopInfo{
			Id:         shopInfoList[i].ShopId,
			Name:       shopInfoList[i].Name,
			AdminName:  shopInfoList[i].AdminName,
			AdminPhone: shopInfoList[i].AdminPhone,
			AdminEmail: shopInfoList[i].AdminEmail,
			Addr:       shopInfoList[i].Addr,
			CreatedAt:  shopInfoList[i].CreatedAt.Format(util.YMD),
			Status:     db.ShopStatusMapping[shopInfoList[i].Status],
		})
	}
	util.EndJson(ctx, rsp)
}

func GetShopListCsvServer(ctx *gin.Context) {
	req := newGetShopListRequest()
	err := ctx.Bind(req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	shopInfoList, _, err := getSortedShopListData(req, true)
	if err != nil {
		logs.CtxError(ctx, "err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	csvRows := make([]shopInfo, 0, len(shopInfoList))
	for _, row := range shopInfoList {
		csvRows = append(csvRows, shopInfo{
			Name:       row.Name,
			AdminName:  row.AdminName,
			AdminPhone: row.AdminPhone,
			AdminEmail: row.AdminEmail,
			Addr:       row.Addr,
			CreatedAt:  row.CreatedAt.Format(util.YMDHMS),
			Status:     db.ShopStatusMapping[row.Status],
		})
	}
	fileName := "用户列表导出.csv"
	file, _ := os.Create(fileName)
	defer file.Close()
	w := csv.NewWriter(file)
	_, _ = file.WriteString("\xEF\xBB\xBF")
	title := make([]string, len(shopInfoColumns))
	for _, item := range shopInfoColumns {
		title = append(title, item.Field)
	}
	_ = w.Write(title)
	for idx := range csvRows {
		var row []string
		var m = make(map[string]string)
		s, _ := json.Marshal(csvRows[idx])
		_ = json.Unmarshal(s, &m)
		for _, k := range shopInfoColumns {
			row = append(row, m[k.Key])
		}
		_ = w.Write(row)
	}
	w.Flush()
	util.SetFileTransportHeader(ctx, fileName)
	_ = os.Remove(fileName)
}

func GetShopListPrintServer(ctx *gin.Context) {
	req := newGetShopListRequest()
	err := ctx.Bind(req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	shopInfoList, totalCnt, err := getSortedShopListData(req, true)
	if err != nil {
		logs.CtxError(ctx, "err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	rsp := getShopListResponse{TotalCnt: totalCnt}
	for i := range shopInfoList {
		rsp.List = append(rsp.List, shopInfo{
			Name:       shopInfoList[i].Name,
			AdminName:  shopInfoList[i].AdminName,
			AdminPhone: shopInfoList[i].AdminPhone,
			AdminEmail: shopInfoList[i].AdminEmail,
			Addr:       shopInfoList[i].Addr,
			CreatedAt:  shopInfoList[i].CreatedAt.Format(util.YMD),
			Status:     db.ShopStatusMapping[shopInfoList[i].Status],
		})
	}
	util.EndJson(ctx, rsp)
}

// 注册新店铺
func AddShopServer(ctx *gin.Context) {
	req := newAddShopRequest()
	err := ctx.Bind(req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	if !isValidShopInfo(req) {
		logs.CtxInfo(ctx, "shop info invalid")
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	email := ctx.GetString("email")
	user, err := db.GetUserByEmail(email)
	if err != nil {
		logs.CtxError(ctx, "get user by email error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	noticeCfg, err := json.Marshal(req.NoticeCfg)
	if err != nil {
		logs.CtxError(ctx, "marshal error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}

	row := &db.ShopList{
		Name:         req.Name,
		UserID:       user.ID,
		Addr:         req.Addr,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		NoticeConfig: string(noticeCfg),
		Status:       0,
		Remark:       req.Remark,
	}
	err = db.SaveShopListRow(row)
	if err != nil {
		logs.CtxError(ctx, "save shop_list row error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	_ = task.SetUserOperationById(user.ID, fmt.Sprintf(util.UserOperation_AddShop, row.ID, row.Name))
	util.EndJson(ctx, nil)
}

// 获取店铺信息
func GetShopInfoServer(ctx *gin.Context) {
	req := &getShopInfoRequest{}
	err := ctx.Bind(req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	rsp, err := getShopInfo(req.ShopId)
	if err != nil {
		logs.CtxError(ctx, "get shop info error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	alarmCnt, err := getShopAlarmCnt(req.ShopId)
	if err != nil {
		logs.CtxError(ctx, "get shop alarm cnt error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	rsp.AlarmCnt = alarmCnt
	util.EndJson(ctx, rsp)
}

// 查询某一用户拥有店铺列表
func GetShopListByEmailServer(ctx *gin.Context) {
	req := &getShopByEmailRequest{}
	err := ctx.Bind(req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	user, err := db.GetUserByEmail(req.UserEmail)
	if err != nil {
		logs.CtxError(ctx, "get user info error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "数据库错误")
		return
	}
	if user == nil {
		logs.CtxInfo(ctx, "user not exist")
		util.ErrorJson(ctx, util.UserNotExist, "用户不存在")
		return
	}
	rows, err := db.GetShopListByUserId(user.ID)
	if err != nil {
		logs.CtxError(ctx, "get shop list by user id error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	rsp := getShopByEmailResponse{}
	for _, row := range rows {
		rsp.List = append(rsp.List, userShopInfo{
			ShopId:   row.ID,
			ShopName: row.Name,
			Addr:     row.Addr,
		})
	}
	util.EndJson(ctx, rsp)
}

func getShopListData(req *getShopListRequest) ([]shopInfoRow, error) {
	begin, err1 := time.Parse(util.YMDHMS, req.CreateBegin)
	end, err2 := time.Parse(util.YMDHMS, req.CreateEnd)
	if err1 != nil || err2 != nil {
		return nil, errors.New("time parse error")
	}
	rows, err := getShopListRows(req.Name, req.Addr, req.Admin, begin, end)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func getSortedShopListData(req *getShopListRequest, all bool) ([]shopInfoRow, int64, error) {
	shopInfoList, err := getShopListData(req)
	if err != nil {
		return nil, 0, err
	}
	if len(req.SortFields) > 0 {
		sort.SliceStable(shopInfoList, func(i, j int) bool {
			a, b := reflect.ValueOf(shopInfoList[i]), reflect.ValueOf(shopInfoList[j])
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
	totalCnt := int64(len(shopInfoList))
	if all {
		return shopInfoList, totalCnt, nil
	}
	offset := req.PageSize * (req.Page - 1)
	if offset+req.PageSize >= totalCnt {
		return shopInfoList[offset:], totalCnt, nil
	}
	return shopInfoList[offset : offset+req.PageSize], totalCnt, nil
}

func newGetShopListRequest() *getShopListRequest {
	return &getShopListRequest{
		Name:        "",
		Admin:       "",
		Addr:        "",
		CreateBegin: time.Now().AddDate(0, 0, -7).Format(util.YMDHMS),
		CreateEnd:   time.Now().Format(util.YMDHMS),
	}
}

func newAddShopRequest() *addShopRequest {
	return &addShopRequest{
		Name:      "",
		Addr:      "",
		Latitude:  0.0,
		Longitude: 0.0,
		NoticeCfg: make(map[string]db.ShopNoticeConfigBase),
		Remark:    "",
	}
}

func isValidShopInfo(req *addShopRequest) bool {
	if req.Name == "" {
		return false
	}
	if req.Addr == "" {
		return false
	}
	return true
}
