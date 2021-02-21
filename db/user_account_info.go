package db

import "gorm.io/gorm"

type UserAccountInfo struct {
	ID       int64  `gorm:"column:id" json:"id"`
	Name     string `gorm:"column:name" json:"name"`
	Email    string `gorm:"column:email" json:"email"`
	Status   int8   `gorm:"column:status" json:"status"`
	Password string `gorm:"column:pwd" json:"password"`
}

func (UserAccountInfo) TableName() string {
	return "user_account_info"
}

func SaveUserInfo(info *UserAccountInfo) (err error) {
	conn, err := fsbmSession.GetConnection()
	if err != nil {
		return err
	}
	err = conn.Save(info).Error
	return
}

func GetUserByEmail(email string) (res *UserAccountInfo, err error) {
	conn, err := fsbmSession.GetConnection()
	if err != nil {
		return
	}
	res = &UserAccountInfo{}
	err = conn.Debug().Where("email = ?", email).First(&res).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func GetAllUser() (res []UserAccountInfo, err error) {
	conn, err := fsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Debug().Order("id").Find(&res).Error
	return
}
