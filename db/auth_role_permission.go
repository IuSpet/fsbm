package db

import "time"

type AuthRolePermission struct {
	ID           int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	RoleID       int64     `gorm:"type:bigint; not null; uniqueIndex:uk_role_permission,priority:1"`
	PermissionID int64     `gorm:"type:bigint; not null; index; uniqueIndex:uk_role_permission,priority:2"`
	Status       int8      `gorm:"type:tinyint; not null; comment:0:正常"`
	CreatedAt    time.Time `gorm:"autoCreateTime; not null"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime; not null"`
}

func (AuthRolePermission) TableName() string {
	return "auth_role_permission"
}

func init() {
	table := AuthRolePermission{}
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
