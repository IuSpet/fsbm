package authority

import (
	"fsbm/db"
	"fsbm/util"
	"time"
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

func getApplyRoleOrderList(user, reviewer string, role []string, status []int8, applyBegin, applyEnd time.Time, reviewBegin, reviewEnd int64) ([]applyRoleRow, error) {
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
	conn = conn.Where("a.created_at between ? and ?", applyBegin, applyEnd)
	if reviewBegin != -1 && reviewEnd != -1 {
		conn = conn.Where("a.review_at >= ? and a.review_at <= ?", reviewBegin, reviewEnd)
	}
	if user != "" {
		conn = conn.Where("b.name like ?", util.LikeCondition(user))
	}
	if len(role) > 0 {
		conn = conn.Where("a.role in ?", role)
	}
	if reviewer != "" {
		conn = conn.Where("c.name like ?", util.LikeCondition(reviewer))
	}
	if len(status) > 0 {
		conn = conn.Where("a.status in ?", status)
	}
	var res []applyRoleRow
	err = conn.Debug().Find(&res).Error
	return res, err
}
