package shop

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"fsbm/db"
	"fsbm/task"
	"fsbm/util"
	"fsbm/util/logs"
	"github.com/gin-gonic/gin"
	"math/rand"
	"os"
	"reflect"
	"sort"
	"strings"
)

var monitorCsvColumns = []struct{ Field, Key string }{
	{"监控名称", "monitor_name"},
	{"店铺名称", "shop_name"},
	{"负责人", "admin_name"},
	{"负责人电话", "admin_phone"},
	{"地址", "addr"},
	{"视频类型", "video_type"},
	{"视频源", "video_src"},
}

// 注册新监控
func AddMonitorServer(ctx *gin.Context) {
	var req addMonitorRequest
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	shopList, err := db.GetShopListById([]int64{req.ShopId})
	if err != nil {
		logs.CtxError(ctx, "get shop list error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	if len(shopList) == 0 {
		logs.CtxError(ctx, "shop not exist")
		util.ErrorJson(ctx, util.ShopNotFound, "店铺信息不存在")
		return
	}
	shop := shopList[0]
	user, err := db.GetUserByEmail(ctx.GetString("email"))
	if err != nil {
		logs.CtxError(ctx, "get user info error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	if user == nil {
		logs.CtxError(ctx, "user not found. email: %s", ctx.GetString("email"))
		util.ErrorJson(ctx, util.AbnormalError, "异常内部错误")
		return
	}
	if shop.UserID != user.ID {
		util.ErrorJson(ctx, util.ShopNotBelongAdjuster, "没有该店铺操作权限")
		return
	}
	row := &db.MonitorList{
		ShopId:    req.ShopId,
		Name:      req.Name,
		VideoType: req.VideoType,
		VideoSrc:  req.VideoSrc,
	}
	err = db.SaveMonitorListRow(row)
	if err != nil {
		logs.CtxError(ctx, "create monitor_list error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	_ = task.SetUserOperationByEmail(ctx.GetString("email"), fmt.Sprintf(util.UserOperation_AddMonitor, row.ShopId, row.ID, row.Name))
	util.EndJson(ctx, nil)
}

// 获取设备列表
func GetMonitorListServer(ctx *gin.Context) {
	req := newGetMonitorListRequest()
	err := ctx.Bind(req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	monitorList, totalCnt, err := getSortedMonitorListData(req, false)
	if err != nil {
		logs.CtxError(ctx, "get device list error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	rsp := getMonitorListResponse{
		List:     monitorList,
		TotalCnt: totalCnt,
	}
	util.EndJson(ctx, rsp)
}

// 获取设备列表csv
func GetMonitorLIstCsvServer(ctx *gin.Context) {
	req := newGetMonitorListRequest()
	err := ctx.Bind(req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	monitorList, _, err := getSortedMonitorListData(req, true)
	if err != nil {
		logs.CtxError(ctx, "get device list error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	fileName := "用户列表导出.csv"
	file, _ := os.Create(fileName)
	defer file.Close()
	w := csv.NewWriter(file)
	_, _ = file.WriteString("\xEF\xBB\xBF")
	title := make([]string, len(shopInfoColumns))
	for _, item := range monitorCsvColumns {
		title = append(title, item.Field)
	}
	_ = w.Write(title)
	for idx := range monitorList {
		var row []string
		var m = make(map[string]string)
		s, _ := json.Marshal(monitorList[idx])
		_ = json.Unmarshal(s, &m)
		for _, k := range monitorCsvColumns {
			row = append(row, m[k.Key])
		}
		_ = w.Write(row)
	}
	w.Flush()
	util.SetFileTransportHeader(ctx, fileName)
	_ = os.Remove(fileName)
}

func GetMonitorListPrintServer(ctx *gin.Context) {
	req := newGetMonitorListRequest()
	err := ctx.Bind(req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	monitorList, totalCnt, err := getSortedMonitorListData(req, true)
	if err != nil {
		logs.CtxError(ctx, "get device list error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	rsp := getMonitorListResponse{
		List:     monitorList,
		TotalCnt: totalCnt,
	}
	util.EndJson(ctx, rsp)
}

// 获取直播墙随机直播源
func GetLiveWallSrcServer(ctx *gin.Context) {
	logs.CtxInfo(ctx, "req live wall src")
	monitorList, err := getLiveMonitorList(16)
	if err != nil {
		logs.CtxError(ctx, "get live src error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	shopIdList := make([]int64, len(monitorList))
	for i, monitor := range monitorList {
		shopIdList[i] = monitor.ShopId
	}
	shopInfoList, err := db.GetShopListById(shopIdList)
	if err != nil {
		logs.CtxError(ctx, "get shop info error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	shopName := make(map[int64]string)
	for _, shop := range shopInfoList {
		shopName[shop.ID] = shop.Name
	}
	rsp := getLiveWallSrcResponse{}
	for _, monitor := range monitorList {
		rsp.List = append(rsp.List, liveSrcInfo{
			MonitorName: monitor.Name,
			ShopName:    shopName[monitor.ShopId],
			VideoType:   monitor.VideoType,
			VideoSrc:    monitor.VideoSrc,
		})
	}
	util.EndJson(ctx, rsp)
}

func getSortedMonitorListData(req *getDeviceListRequest, all bool) ([]monitorInfo, int64, error) {
	rows, err := getMonitorListRows(req.DeviceName, req.ShopName, req.AdminName, req.Addr, req.VideoType)
	if err != nil {
		return nil, 0, err
	}
	if len(req.SortFields) > 0 {
		sort.SliceStable(rows, func(i, j int) bool {
			a, b := reflect.ValueOf(rows[i]), reflect.ValueOf(rows[j])
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
	totalCnt := int64(len(rows))
	if all {
		return rows, totalCnt, nil
	}
	offset := req.PageSize * (req.Page - 1)
	if offset+req.PageSize >= totalCnt {
		return rows[offset:], totalCnt, nil
	}
	return rows[offset : offset+req.PageSize], totalCnt, nil
}

// 从有播放地址的监控中随机选择n个
func getLiveMonitorList(n int) ([]db.MonitorList, error) {
	rows, err := db.GetLiveMonitorRows()
	if err != nil {
		return nil, err
	}
	rand.Shuffle(len(rows), func(i, j int) {
		rows[i], rows[j] = rows[j], rows[i]
	})
	if len(rows) > n {
		return rows[:n], nil
	}
	return rows, nil
}

func newGetMonitorListRequest() *getDeviceListRequest {
	return &getDeviceListRequest{
		DeviceName: "",
		ShopName:   "",
		AdminName:  "",
		Addr:       "",
		VideoType:  "",
		ListReqField: util.ListReqField{
			Page:       1,
			PageSize:   10,
			SortFields: nil,
		},
	}
}
