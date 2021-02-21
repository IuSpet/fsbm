package db

type AuthRole struct {
	ID     int64  `gorm:"column:id" json:"id"`
	Role   string `gorm:"column:role" json:"role"`
	Type   string `gorm:"column:type" json:"type"`
	Status int8   `gorm:"column:status" json:"status"`
}

func (AuthRole) TableName() string {
	return "auth_role"
}

func GetRoleById(id int64) (res []AuthRole, err error) {
	sqlFmt := `
	select * from auth_role a join auth_user_role b on a.id = b.role_id where b.id = ?
`
	conn, err := fsbmSession.GetConnection()
	if err != nil {
		return nil, err
	}
	err = conn.Debug().Raw(sqlFmt, id).Find(&res).Error
	return
}

func SaveRole(role *AuthRole) (err error) {
	conn, err := fsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Save(role).Error
	return
}

func GetRoleByName(t string, name string) (res *AuthRole, err error) {
	conn, err := fsbmSession.GetConnection()
	if err != nil {
		return
	}
	res = &AuthRole{}
	err = conn.Debug().Where("type = ? and role = ?", t, name).Find(res).Error
	return
}
