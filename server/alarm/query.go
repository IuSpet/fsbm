package alarm

import (
	"fsbm/db"
	"fsbm/util"
)

func queryAlarmList(shopName, adminName, startTime, endTime string, alarmType []int8) (res []alarmListRow, err error) {
	conn, err := db.FsbmSession.GetConnection()
	if err != nil {
		return
	}
	conn = conn.Table("record_alarm a " +
		"LEFT JOIN shop_list b ON a.shop_id = b.id " +
		"LEFT JOIN user_account_info c ON a.user_id = c.id")
	conn = conn.Select("b.name as shop_name," +
		"c.name as admin_name," +
		"c.phone as admin_phone," +
		"b.addr," +
		"a.alarm_type," +
		"a.id as alarm_id," +
		"a.alarm_at")
	if shopName != "" {
		conn = conn.Where("shop_name like ?", util.LikeCondition(shopName))
	}
	if adminName != "" {
		conn = conn.Where("admin_name like ", util.LikeCondition(adminName))
	}
	if startTime != "" && endTime != "" {
		conn = conn.Where("alarm_at >= ? and alarm_at <= ?", startTime, endTime)
	}
	if len(alarmType) > 0 {
		conn = conn.Where("alarm_type in ?", alarmType)
	}
	err = conn.Debug().Find(&res).Error
	return
}

func getAlarmInfo(id int64) (res alarmDetailInfo, err error) {
	conn, err := db.FsbmSession.GetConnection()
	if err != nil {
		return
	}
	conn = conn.Table("record_alarm a " +
		"left join shop_list b on a.shop_id = b.id " +
		"left join user_account_info c on a.user_id = c.id")
	conn = conn.Select("b.name as shop_name," +
		"b.addr as addr," +
		"c.name as admin_name," +
		"c.phone as admin_phone," +
		"c.email as admin_email," +
		"a.alarm_type," +
		"a.alarm_at")
	conn = conn.Where("a.id = ?", id)
	err = conn.Debug().Find(&res).Error
	return
}
