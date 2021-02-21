package db

type AuthRolePermission struct {
	ID           int64 `gorm:"column:id" json:"id"`
	RoleID       int64 `gorm:"column:role_id" json:"role_id"`
	PermissionID int64 `gorm:"column:permission_id" json:"permission_id"`
}

func (AuthRolePermission) TableName() string {
	return "auth_role_permission"
}
