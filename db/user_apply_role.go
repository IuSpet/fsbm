package db

import "time"

type UserApplyRole struct {
	ID        int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	UserID    int64     `gorm:"type:bigint; not null; uniqueIndex:uk_user_role,priority:1"`
	RoleID    int64     `gorm:"type:bigint; not null; index; uniqueIndex:uk_user_role,priority:2"`
	Status    int8      `gorm:"type:tinyint; not null; comment:1:申请中,2:已通过,3:未通过"`
	CreatedAt time.Time `gorm:"autoCreateTime; not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime; not null"`
}

func (UserApplyRole) TableName() string {
	return "user_apply_role"
}

func init(){
	table := UserApplyRole{}
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

func SaveUserApplyRoleRows(rows []UserApplyRole) (err error) {
	conn, err := fsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Debug().Save(rows).Error
	return
}

func GetUserApplyRows(userID int64) (res []UserApplyRole, err error) {
	conn, err := fsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Where("user_id = ? and status = 1", userID).Find(&res).Error
	return
}
