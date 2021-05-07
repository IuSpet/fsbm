package manager

import "fsbm/db"

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
