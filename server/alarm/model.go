package alarm

import (
	"fsbm/util"
)

type alarmListRequest struct {
	ShopName  string `json:"shop_name"`
	AdminName string `json:"admin_name"`
	AlarmType []int8 `json:"alarm_type"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	util.ListReqField
}

type alarmListResponse struct {
	List     []alarmInfo `json:"list"`
	TotalCnt int64       `json:"total_cnt"`
}

type alarmInfo struct {
	AlarmId      int64  `json:"alarm_id"`
	ShopName     string `json:"shop_name"`
	AdminName    string `json:"admin_name"`
	AdminPhone   string `json:"admin_phone"`
	Addr         string `json:"addr"`
	AlarmContent string `json:"alarm_content"`
	AlarmAt      string `json:"alarm_at"`
}

type alarmListRow struct {
	AlarmId    int64
	ShopName   string
	AdminName  string
	AdminPhone string
	Addr       string
	AlarmType  int8
	AlarmAt    string
}

type alarmListCsvRow struct {
	ShopName     string `json:"shop_name"`
	AdminName    string `json:"admin_name"`
	AdminPhone   string `json:"admin_phone"`
	Addr         string `json:"addr"`
	AlarmContent string `json:"alarm_content"`
	AlarmAt      string `json:"alarm_at"`
	Detail       string `json:"detail"`
}

type alarmDetailInfoRequest struct {
	AlarmId int64 `json:"alarm_id"`
}

type alarmDetailInfoResponse struct {
	Info alarmDetailInfo `json:"info"`
}

type alarmDetailInfo struct {
	ShopName     string `json:"shop_name"`
	Addr         string `json:"addr"`
	AdminName    string `json:"admin_name"`
	AdminPhone   string `json:"admin_phone"`
	AdminEmail   string `json:"admin_email"`
	AlarmContent string `json:"alarm_content"`
	AlarmType    int8   `json:"alarm_type"`
	AlarmAt      string `json:"alarm_at"`
}
