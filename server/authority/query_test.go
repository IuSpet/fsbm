package authority

import (
	"fmt"
	"fsbm/conf"
	"fsbm/db"
	"testing"
)

func TestGetUserRoleList(t *testing.T) {
	conf.Init()
	db.Init()
	res, err := getUserRoleList(1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}

func TestGetApplyRoleOrderList(t *testing.T) {
	_, err := getApplyRoleOrderList("user", "role", "reviewer", []int8{0, 1}, 0, 123120388)
	if err != nil {
		t.Error(err)
	}
}
