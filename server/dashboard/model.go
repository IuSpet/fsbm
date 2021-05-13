package dashboard

type dashboardStatsResponse struct {
	RecordCnt    int64   `json:"record_cnt"`
	AlarmCnt     int64   `json:"alarm_cnt"`
	LatestRecord string  `json:"latest_record"`
	ShopPassRate float64 `json:"shop_pass_rate"`
}
