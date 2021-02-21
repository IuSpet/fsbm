package db

type UserApplyRole struct {
	ID     int64 `gorm:"column:id" json:"id"`
	UserID int64 `gorm:"column:user_id" json:"user_id"`
	RoleID int64 `gorm:"column:role_id" json:"role_id"`
	Status int8  `gorm:"column:status" json:"status"`
}

func (UserApplyRole) TableName() string {
	return "user_apply_role"
}

func SaveUserApplyRoleRows(rows []UserApplyRole) (err error) {
	conn, err := fsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Debug().Save(rows).Error
	return
}

func GetUserApplyRows(userID int64) (res []UserApplyRole, err error) {
	conn, err := fsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Where("user_id = ? and status = 1",userID).Find(&res).Error
	return
}
