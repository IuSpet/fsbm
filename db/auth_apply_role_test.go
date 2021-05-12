package db

import (
	"fmt"
	"fsbm/conf"
	"testing"
)

func TestGetAuthApplyRoleById(t *testing.T) {
	conf.Init()
	Init()
	res, err := GetAuthApplyRoleById(500)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(err)
	fmt.Printf("%+v", res)
}
