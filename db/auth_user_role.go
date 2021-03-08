package db

import "time"

type AuthUserRole struct {
	ID        int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	UserID    int64     `gorm:"type:bigint; not null; uniqueIndex:uk_user_role,priority:1"`
	RoleID    int64     `gorm:"type:bigint; not null; index; uniqueIndex:uk_user_role,priority:2"`
	Status    int8      `gorm:"type:tinyint; not null; comment:0:正常"`
	CreatedAt time.Time `gorm:"autoCreateTime; not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime; not null"`
}

func (AuthUserRole) TableName() string {
	return "auth_user_role"
}

func init() {
	table := AuthUserRole{}
	RegisterMigration(table.TableName(), func() {
		conn, err := fsbmSession.GetConnection()
		if err != nil {
			panic(err)
		}
		err = conn.Set("gorm:table_options", "ENGINE=INNODB CHARSET=utf8").AutoMigrate(&table)
		if err != nil {
			panic(err)
		}
	})
}

func SaveAuthUserRoleRows(rows []AuthUserRole) (err error) {
	conn, err := fsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Save(rows).Error
	return
}

func RemoveUserRole(userID int64, roleIDList []int64) (err error) {
	conn, err := fsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Where("user_id = ? and role_id in (?)", userID, roleIDList).Delete(AuthUserRole{}).Error
	return
}
