package shop

import (
	"fsbm/db"
	"fsbm/util"
	"time"
)

func getShopListRows(name, addr, admin string, begin, end time.Time) (rows []shopInfoRow, err error) {
	conn, err := db.FsbmSession.GetConnection()
	if err != nil {
		return
	}
	conn = conn.Select("a.name as name," +
		"a.addr as addr," +
		"a.created_at as created_at," +
		"a.status as status," +
		"b.name as admin_name," +
		"b.phone as admin_phone," +
		"b.email as admin_email," +
		"a.id as shop_id")
	conn = conn.Table("shop_list a left join user_account_info b on a.user_id = b.id")
	if name != "" {
		conn = conn.Where("a.name like ?", util.LikeCondition(name))
	}
	if addr != "" {
		conn = conn.Where("a.addr like ?", util.LikeCondition(addr))
	}
	if admin != "" {
		conn = conn.Where("b.name = ?", admin)
	}
	conn = conn.Where("a.created_at between ? and ?", begin, end)
	err = conn.Debug().Find(&rows).Error
	return
}

func getMonitorListRows(name, shop, admin, addr, videoType string) (rows []monitorInfo, err error) {
	conn, err := db.FsbmSession.GetConnection()
	if err != nil {
		return
	}
	conn = conn.Select("a.name as monitor_name," +
		"a.video_type as video_type," +
		"a.video_src as video_src," +
		"b.name as shop_name," +
		"b.addr as addr," +
		"c.name as user_name," +
		"c.phone as user_phone")
	conn = conn.Table("monitor_list a " +
		"LEFT JOIN shop_list b ON a.shop_id = b.id " +
		"LEFT JOIN user_account_info c ON b.user_id = c.id")
	if name != "" {
		conn = conn.Where("a.name like ?", util.LikeCondition(name))
	}
	if shop != "" {
		conn = conn.Where("b.name like ?", util.LikeCondition(shop))
	}
	if admin != "" {
		conn = conn.Where("c.name like ?", util.LikeCondition(admin))
	}
	if addr != "" {
		conn = conn.Where("b.addr like ?", util.LikeCondition(addr))
	}
	if videoType != "" {
		conn = conn.Where("a.video_type = ?", videoType)
	}
	err = conn.Debug().Find(&rows).Error
	return
}

func getShopInfo(shopId int64) (res getShopInfoResponse, err error) {
	conn, err := db.FsbmSession.GetConnection()
	if err != nil {
		return
	}
	conn = conn.Select("a.name as shop_name," +
		"a.addr," +
		"b.name as user_name," +
		"b.email as user_email," +
		"b.phone as user_phone")
	conn = conn.Table("shop_list a left join user_account_info b on a.user_id = b.id")
	conn = conn.Where("a.id = ?", shopId)
	err = conn.Find(&res).Error
	return
}

func getShopAlarmCnt(shopId int64) (alarmCnt int64, err error) {
	conn, err := db.FsbmSession.GetConnection()
	if err != nil {
		return
	}
	var rows []db.RecordAlarm
	err = conn.Where("shop_id = ?", shopId).Find(&rows).Error
	return int64(len(rows)), err
}
