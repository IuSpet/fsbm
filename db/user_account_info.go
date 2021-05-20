package db

import (
	"gorm.io/gorm"
	"time"
)

type UserAccountInfo struct {
	ID        int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	Name      string    `gorm:"type:varchar(127); not null "`
	Email     string    `gorm:"type:varchar(127); not null; uniqueIndex"`
	Status    int8      `gorm:"type:tinyint; not null; comment:0:正常,1:已删除"`
	Password  string    `gorm:"type:varchar(127); not null"`
	Phone     string    `gorm:"type:varchar(127); not null"`
	Gender    int8      `gorm:"type:tinyint; not null; default:0; comment:0:未设置,1:男,2:女"`
	Age       int8      `gorm:"type:tinyint; not null; default:0"`
	Avatar    []byte    `gorm:"type:blob"`
	CreatedAt time.Time `gorm:"autoCreateTime; not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime; not null"`
}

func (UserAccountInfo) TableName() string {
	return "user_account_info"
}

var UserStatusMapping = map[int8]string{
	0: "正常",
	1: "已删除",
}

var UserGenderMapping = map[int8]string{
	0: "其他",
	1: "男",
	2: "女",
}

func init() {
	table := UserAccountInfo{}
	RegisterMigration(table.TableName(), func() {
		conn, err := FsbmSession.GetConnection()
		if err != nil {
			panic(err)
		}
		err = conn.Set("gorm:table_options", "ENGINE=INNODB CHARSET=utf8").AutoMigrate(&table)
		if err != nil {
			panic(err)
		}
	})
}

func SaveUserInfo(info *UserAccountInfo) (err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return err
	}
	err = conn.Save(info).Error
	return
}

func GetUserByEmail(email string) (res *UserAccountInfo, err error) {
	conn, err := FsbmSession.GetConnection()
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

func GetUserById(id int64) (res *UserAccountInfo, err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	res = &UserAccountInfo{ID: id}
	err = conn.First(res).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func GetAllUser() (res []UserAccountInfo, err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Debug().Order("id").Find(&res).Error
	return
}

func FuzzySearchUser(conditions []string) (res []UserAccountInfo, err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	for _, condition := range conditions {
		conn.Where(condition)
	}
	err = conn.Find(&res).Error
	return
}

func SetAvatar(email string, avatar []byte) (err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Model(&UserAccountInfo{}).Where("email = ?", email).Update("avatar", avatar).Error
	return
}

func GetUserAccountInfoTotalCnt() (res int64, err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Count(&res).Error
	return
}

func GetUserAccountInfoById(id int64) (res *UserAccountInfo, err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	res = &UserAccountInfo{ID: id}
	err = conn.First(res).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}
