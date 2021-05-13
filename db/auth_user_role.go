package db

import (
	"gorm.io/gorm"
	"time"
)

const (
	AuthUserRoleStatus_Active  int8 = 0
	AuthUserRoleStatus_Expired int8 = 1
)

var authUserRoleStatusMapping = map[int8]string{
	AuthUserRoleStatus_Active:  "正常",
	AuthUserRoleStatus_Expired: "已过期",
}

// 用户id与角色id表中唯一，已有关系记录只修改生效时间与状态
type AuthUserRole struct {
	ID        int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	UserID    int64     `gorm:"type:bigint; not null; uniqueIndex:uk_user_role,priority:1"`
	RoleID    int64     `gorm:"type:bigint; not null; index; uniqueIndex:uk_user_role,priority:2"`
	StartTime time.Time `gorm:"not null; comment: 开始时间"`
	EndTime   time.Time `gorm:"not null; comment: 结束时间"`
	Status    int8      `gorm:"type:tinyint; not null; comment:0:正常,1:已过期"`
	CreatedAt time.Time `gorm:"autoCreateTime; not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime; not null"`
}

func (AuthUserRole) TableName() string {
	return "auth_user_role"
}

func init() {
	table := AuthUserRole{}
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

func SaveAuthUserRoleRow(row *AuthUserRole) (err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Save(row).Error
	return
}

func SaveAuthUserRoleRows(rows []AuthUserRole) (err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Save(rows).Error
	return
}

func RemoveUserRole(userID int64, roleIDList []int64) (err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Where("user_id = ? and role_id in (?)", userID, roleIDList).Delete(AuthUserRole{}).Error
	return
}

// 获取用户激活角色
func GetUserActiveRoles(userId int64) (res []AuthUserRole, err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Where("user_id = ? and status = 0", userId).Find(&res).Error
	return
}

// 获取用户过期角色
func GetUserExpiredRoles(userId int64) (res []AuthUserRole, err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Where("user_id = ? and status = 1", userId).Find(&res).Error
	return
}

// 获取用户角色关联
func GetUserRoleRow(userId, RoleId int64) (res *AuthUserRole, err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Where("user_id = ? and role_id = 1", userId, RoleId).First(&res).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

// 获取所有激活关系
func GetAllActiveRelation() (res []AuthUserRole, err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Where("status = 0").Find(&res).Error
	return
}
