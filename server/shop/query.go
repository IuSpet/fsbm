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
		"b.email as admin_email")
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
