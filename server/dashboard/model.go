package dashboard

type globalStatsResponse struct {
	RecordCnt    int64   `json:"record_cnt"`
	AlarmCnt     int64   `json:"alarm_cnt"`
	LatestRecord string  `json:"latest_record"`
	ShopPassRate float64 `json:"shop_pass_rate"`
}

type mapShopInfoListResponse struct {
	List []mapShopInfo `json:"list"`
}

type mapShopInfo struct {
	ShopId    int64   `json:"shop_id"`
	ShopName  string  `json:"shop_name"`
	UserId    int64   `json:"user_id"`
	UserName  string  `json:"user_name"`
	UserPhone string  `json:"user_phone"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	AlarmCnt  int64   `json:"alarm_cnt"`
}

type shopAlarmCnt struct {
	ShopId int64
	Cnt    int64
}
