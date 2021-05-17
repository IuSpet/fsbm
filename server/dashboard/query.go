package dashboard

import (
	"fsbm/db"
	"time"
)

func getTodayShopAlarm() (res []shopAlarmCnt, err error) {
	conn, err := db.FsbmSession.GetConnection()
	if err != nil {
		return
	}
	now := time.Now()
	beg := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	conn = conn.Select("shop_id,count(1) as cnt")
	conn = conn.Table("record_alarm")
	conn = conn.Where("alarm_at >= ?", beg)
	conn = conn.Group("shop_id")
	err = conn.Debug().Find(&res).Error
	return
}

func getShopInfoList() (res []mapShopInfo, err error) {
	conn, err := db.FsbmSession.GetConnection()
	if err != nil {
		return
	}
	conn = conn.Select("a.id as shop_id," +
		"a.name as shop_name, " +
		"b.id as user_id, " +
		"b.name as user_name," +
		"b.phone as user_phone," +
		"latitude," +
		"longitude ")
	conn = conn.Table("shop_list a " +
		"join user_account_info b on a.user_id = b.id")
	conn = conn.Where("a.status = 0")
	err = conn.Find(&res).Error
	return
}
