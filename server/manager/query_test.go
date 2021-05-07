package manager

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
