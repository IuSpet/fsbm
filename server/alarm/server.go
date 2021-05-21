package alarm

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
	"strings"
	"time"
)

var alarmInfoColumns = []struct{ Field, Key string }{
	{"店铺名称", "shop_name"},
	{"负责人名称", "admin_name"},
	{"负责人电话", "admin_phone"},
	{"店铺地址", "addr"},
	{"报警内容", "alarm_content"},
	{"报警时间", "alarm_at"},
	{"详情", "detail"},
}

// 报警列表接口
func AlarmListServer(ctx *gin.Context) {
	req := newAlarmListRequest()
	err := ctx.Bind(req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	alarmList, totalCnt, err := getSortedAlarmList(req, false)
	if err != nil {
		logs.CtxError(ctx, "get alarm list error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	rsp := alarmListResponse{TotalCnt: totalCnt}
	for _, row := range alarmList {
		rsp.List = append(rsp.List, alarmInfo{
			AlarmId:      row.AlarmId,
			ShopName:     row.ShopName,
			AdminName:    row.AdminName,
			AdminPhone:   row.AdminPhone,
			Addr:         row.Addr,
			AlarmContent: db.RecordAlarmAlarmTypeMapping[row.AlarmType],
			AlarmAt:      row.AlarmAt,
		})
	}
	util.EndJson(ctx, rsp)
}

// 报警列表csv接口
func AlarmListCsvServer(ctx *gin.Context) {
	req := newAlarmListRequest()
	err := ctx.Bind(req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	alarmList, _, err := getSortedAlarmList(req, true)
	if err != nil {
		logs.CtxError(ctx, "get alarm list error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	csvRows := make([]alarmListCsvRow, 0, len(alarmList))
	for _, row := range alarmList {
		csvRows = append(csvRows, alarmListCsvRow{
			ShopName:     row.ShopName,
			AdminName:    row.AdminName,
			AdminPhone:   row.AdminPhone,
			Addr:         row.Addr,
			AlarmContent: db.RecordAlarmAlarmTypeMapping[row.AlarmType],
			AlarmAt:      row.AlarmAt,
			Detail:       fmt.Sprintf("http://47.95.248.242/#/alarm/alarm_detail?id=%d", row.AlarmId),
		})
	}
	fileName := "用户列表导出.csv"
	file, _ := os.Create(fileName)
	defer file.Close()
	w := csv.NewWriter(file)
	_, _ = file.WriteString("\xEF\xBB\xBF")
	title := make([]string, len(alarmInfoColumns))
	for _, item := range alarmInfoColumns {
		title = append(title, item.Field)
	}
	_ = w.Write(title)
	for idx := range csvRows {
		var row []string
		var m = make(map[string]string)
		s, _ := json.Marshal(csvRows[idx])
		_ = json.Unmarshal(s, &m)
		for _, k := range alarmInfoColumns {
			row = append(row, m[k.Key])
		}
		_ = w.Write(row)
	}
	w.Flush()
	util.SetFileTransportHeader(ctx, fileName)
	_ = os.Remove(fileName)
}

func AlarmListPrintServer(ctx *gin.Context) {
	req := newAlarmListRequest()
	err := ctx.Bind(req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	alarmList, totalCnt, err := getSortedAlarmList(req, true)
	if err != nil {
		logs.CtxError(ctx, "get alarm list error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	rsp := alarmListResponse{TotalCnt: totalCnt}
	for _, row := range alarmList {
		rsp.List = append(rsp.List, alarmInfo{
			AlarmId:      row.AlarmId,
			ShopName:     row.ShopName,
			AdminName:    row.AdminName,
			AdminPhone:   row.AdminPhone,
			Addr:         row.Addr,
			AlarmContent: db.RecordAlarmAlarmTypeMapping[row.AlarmType],
			AlarmAt:      row.AlarmAt,
		})
	}
	util.EndJson(ctx, rsp)
}

func AlarmDetailInfoServer(ctx *gin.Context) {
	req := &alarmDetailInfoRequest{AlarmId: -1}
	err := ctx.Bind(req)
	if err != nil {
		logs.CtxError(ctx, "bind param error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	info, err := getAlarmInfo(req.AlarmId)
	if err != nil {
		logs.CtxError(ctx, "get alarm info error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	info.AlarmContent = db.RecordAlarmAlarmTypeMapping[info.AlarmType]
	rsp := alarmDetailInfoResponse{Info: info}
	util.EndJson(ctx, rsp)
}

func getSortedAlarmList(req *alarmListRequest, all bool) ([]alarmListRow, int64, error) {
	alarmList, err := queryAlarmList(req.ShopName, req.AdminName, req.StartTime, req.EndTime, req.AlarmType)
	if err != nil {
		return nil, 0, err
	}
	if len(req.SortFields) > 0 {
		sort.SliceStable(alarmList, func(i, j int) bool {
			a, b := reflect.ValueOf(alarmList[i]), reflect.ValueOf(alarmList[j])
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
	totalCnt := int64(len(alarmList))
	if all {
		return alarmList, totalCnt, nil
	}
	offset := req.PageSize * (req.Page - 1)
	if offset+req.PageSize >= totalCnt {
		return alarmList[offset:], totalCnt, nil
	}
	return alarmList[offset : offset+req.PageSize], totalCnt, nil
}

func newAlarmListRequest() *alarmListRequest {
	return &alarmListRequest{
		ShopName:  "",
		AdminName: "",
		AlarmType: nil,
		StartTime: time.Unix(0, 0).Format(util.YMDHMS),
		EndTime:   time.Now().Format(util.YMDHMS),
		ListReqField: util.ListReqField{
			Page:       1,
			PageSize:   10,
			SortFields: nil,
		},
	}
}
