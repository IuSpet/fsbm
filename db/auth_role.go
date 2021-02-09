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

