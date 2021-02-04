package db

import "time"

type UserAccountInfo struct {
	Id         int64     `gorm:"column:id" json:"id"`
	Name       string    `gorm:"column:name" json:"name"`
	Email      string    `gorm:"column:email" json:"email"`
	Status     int8      `gorm:"column:status" json:"status"`
	Password   string    `gorm:"column:pwd" json:"password"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	ModifyTime time.Time `gorm:"column:modify_time" json:"modify_time"`
}

func (UserAccountInfo) TableName() string {
	return "user_account_info"
}

func SaveUserInfo(info *UserAccountInfo) error {
	conn, err := fsbmSession.GetConnection()
	if err != nil {
		return err
	}
	conn.Save(info)
	return nil
}
