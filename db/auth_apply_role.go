package db

import "time"

var AuthApplyRoleStatusMapping = map[int8]string{
	0: "未审核",
	1: "通过",
	2: "拒绝",
}

const (
	AuthApplyRoleStatus_Unreviewd int8 = 0
	AuthApplyRoleStatus_Approve   int8 = 1
	AuthApplyRoleStatus_Deny      int8 = 2
)

// 角色申请工单表
type AuthApplyRole struct {
	ID           int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	UserId       int64     `gorm:"not null; comment:用户id"`
	Email        string    `gorm:"type:varchar(127);not null; comment:用户邮箱"`
	RoleId       int64     `gorm:"not null; comment:角色id"`
	Role         string    `gorm:"type:varchar(127);not null; comment:角色"`
	Reason       string    `gorm:"type:text;not null; comment:申请理由"`
	ReviewUserId int64     `gorm:"not null; comment:审核用户id"`
	ReviewTs     int64     `gorm:"not null; comment:审核时间"`
	ReviewReason string    `gorm:"type:text;not null; comment:审核理由"`
	Status       int8      `gorm:"type:tinyint; not null; comment:0:未审核,1:通过,2:拒绝"`
	CreatedAt    time.Time `gorm:"autoCreateTime; not null"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime; not null"`
}

func (AuthApplyRole) TableName() string {
	return "auth_permission"
}

func init() {
	table := AuthApplyRole{}
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

func SaveAthApplyRoleRow(row *AuthApplyRole) (err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Save(row).Error
	return
}
