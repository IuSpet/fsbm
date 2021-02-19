package db

type AuthPermission struct {
	ID         int64  `gorm:"column:id" json:"id"`
	Permission string `gorm:"column:permission" json:"permission"`
	Type       string `gorm:"column:type" json:"type"`
	Status     int8   `gorm:"column:status" json:"status"`
}

func (AuthPermission) TableName() string {
	return "auth_permission"
}

func GetPermissionByRoleID(roleID []int64) (res []AuthPermission, err error) {
	sqlFmt := `
	select * from auth_permission a join auth_role_permission b on a.id = b.permission_id where b.role_id in (?)
`
	conn, err := fsbmSession.GetConnection()
	if err != nil {
		return nil, err
	}
	err = conn.Raw(sqlFmt, roleID).Find(&res).Error
	return
}

func (a AuthPermission) String() string {
	return a.Type + ":" + a.Permission
}
