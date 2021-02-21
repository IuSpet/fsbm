package db

import (
	"fmt"
	"testing"
)

func TestSaveUserApplyRoleRows(t *testing.T) {
	row := UserApplyRole{
		UserID: 1,
		RoleID: 2,
		Status: 1,
	}
	err := SaveUserApplyRoleRows([]UserApplyRole{row})
	if err != nil {
		panic(err)
	}
}

func TestGetUserApplyRows(t *testing.T) {
	res, err := GetUserApplyRows(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", res)
}
