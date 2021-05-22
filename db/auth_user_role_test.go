package db

import (
	"fsbm/conf"
	"testing"
	"time"
)

func TestRemoveUserRole(t *testing.T) {
	err := RemoveUserRole(5, []int64{1, 2})
	if err != nil {
		panic(err)
	}
}
func TestSaveAuthUserRoleRow(t *testing.T) {
	conf.Init()
	Init()
	err := SaveAuthUserRoleRow(&AuthUserRole{
		UserID:    6,
		RoleID:    1,
		StartTime: time.Now(),
		EndTime:   time.Now().AddDate(100, 0, 0),
		Status:    AuthUserRoleStatus_Active,
	})
	if err != nil {
		panic(err)
	}
}
