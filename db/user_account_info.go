package db

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
	err = conn.Debug().Where("email = ? and status = 0", email).First(&res).Error
	return
}
