package db

import "time"

type AuthPermission struct {
	ID         int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	Permission string    `gorm:"type:varchar(127); not null"`
	Type       string    `gorm:"type:varchar(127); not null"`
	Status     int8      `gorm:"type:tinyint; not null; comment:0:正常"`
	CreatedAt  time.Time `gorm:"autoCreateTime; not null"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime; not null"`
}

func (AuthPermission) TableName() string {
	return "auth_permission"
}

func init() {
	table := AuthPermission{}
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

func GetPermissionByRoleID(roleID []int64) (res []AuthPermission, err error) {
	sqlFmt := `
	select * from auth_permission a join auth_role_permission b on a.id = b.permission_id where b.role_id in (?)
`
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return nil, err
	}
	err = conn.Raw(sqlFmt, roleID).Find(&res).Error
	return
}

func (a AuthPermission) String() string {
	return a.Type + ":" + a.Permission
}
