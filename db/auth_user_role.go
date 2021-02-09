package db

type AuthUserRole struct {
	ID     int64 `gorm:"column:id" json:"id"`
	UserID int64 `gorm:"column:user_id" json:"user_id"`
	RoleID int64 `gorm:"column:role_id" json:"role_id"`
	Status int8  `gorm:"column:status" json:"status"`
}

func (AuthUserRole) TableName() string {
	return "auth_user_role"
}
