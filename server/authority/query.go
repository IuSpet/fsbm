package authority

import (
	"fsbm/db"
)

func getUserRoleList(userId int64) ([]userRoleStatusInfo, error) {
	conn, err := db.FsbmSession.GetConnection()
	if err != nil {
		return nil, err
	}
	sqlFmt := `select a.status,a.role_id,b.role from auth_user_role a left join auth_role b on a.role_id = b.id where a.user_id = ?`
	var res []userRoleStatusInfo
	err = conn.Raw(sqlFmt, userId).Find(&res).Error
	return res, err
}

func getApplyRoleOrderList(user, role, reviewer string, status []int8, applyBegin, applyEnd, reviewBegin, reviewEnd int64) ([]applyRoleRow, error) {
	conn, err := db.FsbmSession.GetConnection()
	if err != nil {
		return nil, err
	}
	conn = conn.Select("a.id," +
		"b.name as user," +
		"a.role," +
		"a.reason," +
		"a.status," +
		"c.name AS reviewer," +
		"a.review_reason," +
		"a.review_at," +
		"a.created_at")
	conn = conn.Table("auth_apply_role a " +
		"LEFT JOIN user_account_info b ON a.user_id = b.id " +
		"LEFT JOIN user_account_info c ON a.review_user_id = c.id ")
	conn = conn.Where("created_at between ? and ?", applyBegin, applyEnd)
	conn = conn.Where("review_at >= ? and review_at <= ?", reviewBegin, reviewEnd)
	if user != "" {
		conn = conn.Where("b.name = ?", user)
	}
	if role != "" {
		conn = conn.Where("a.role = ?", role)
	}
	if reviewer != "" {
		conn = conn.Where("c.name = ?", reviewer)
	}
	if status != nil {
		conn = conn.Where("a.status in ?", status)
	}
	var res []applyRoleRow
	err = conn.Debug().Find(&res).Error
	return res, err
}
