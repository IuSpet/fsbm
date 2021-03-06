package db

import (
	"gorm.io/gorm"
	"time"
)

type UserAccountInfo struct {
	ID        int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	Name      string    `gorm:"type:varchar(128); not null "`
	Email     string    `gorm:"type:varchar(128); not null; uniqueIndex"`
	Status    int8      `gorm:"type:tinyint; not null; comment:0:正常,1:已删除"`
	Password  string    `gorm:"type:varchar(128); not null"`
	Phone     string    `gorm:"type:bigint; not null"`
	CreatedAt time.Time `gorm:"autoCreateTime; not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime; not null"`
}

func (UserAccountInfo) TableName() string {
	return "user_account_info"
}

func init() {
	table := UserAccountInfo{}
	RegisterMigration(table.TableName(), func() {
		conn, err := fsbmSession.GetConnection()
		if err != nil {
			panic(err)
		}
		err = conn.Debug().Set("gorm:table_options", "ENGINE=INNODB CHARSET=utf8").AutoMigrate(&table)
		if err != nil {
			panic(err)
		}
	})
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
	err = conn.Debug().Where("email = ?", email).First(res).Error
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

func FuzzySearchUser(conditions []string) (res []UserAccountInfo, err error) {
	conn, err := fsbmSession.GetConnection()
	if err != nil {
		return
	}
	for _, condition := range conditions {
		conn.Where(condition)
	}
	err = conn.Find(&res).Error
	return
}
