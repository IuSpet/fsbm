package dashboard

import (
	"fsbm/db"
	"time"
)

func getTodayShopAlarm() (res []db.RecordAlarm, err error) {
	conn, err := db.FsbmSession.GetConnection()
	if err != nil {
		return
	}
	now := time.Now()
	beg := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	err = conn.Debug().Where("alarm_at >= ?", beg).Group("shop_id").Find(&res).Error
	return
}
