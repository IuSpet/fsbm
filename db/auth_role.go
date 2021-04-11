package db

import "time"

type AuthRole struct {
	ID        int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	Role      string    `gorm:"type:varchar(128); not null; index; uniqueIndex:uk_type_role,priority:2"`
	Type      string    `gorm:"type:varchar(128); not null; uniqueIndex:uk_type_role,priority:1"`
	Status    int8      `gorm:"type:tinyint; not null; comment:0:正常"`
	CreatedAt time.Time `gorm:"autoCreateTime; not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime; not null"`
}

func (AuthRole) TableName() string {
	return "auth_role"
}

func init() {
	table := AuthRole{}
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

func GetRoleById(id int64) (res []AuthRole, err error) {
	sqlFmt := `
	select * from auth_role a join auth_user_role b on a.id = b.role_id where b.id = ?
`
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return nil, err
	}
	err = conn.Debug().Raw(sqlFmt, id).Find(&res).Error
	return
}

func SaveRole(role *AuthRole) (err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Save(role).Error
	return
}

func GetRoleByName(t string, name string) (res *AuthRole, err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	res = &AuthRole{}
	err = conn.Debug().Where("type = ? and role = ?", t, name).Find(res).Error
	return
}
