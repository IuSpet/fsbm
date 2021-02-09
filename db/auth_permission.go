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
