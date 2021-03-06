package db

import (
	"fmt"
	"fsbm/conf"
	"testing"
)

func TestAddRole(t *testing.T) {
	conf.Init()
	Init()
	role := &AuthRole{
		Role:   "normal_user",
		Type:   "user",
		Status: 0,
	}
	err := SaveRole(role)
	if err != nil {
		panic(err)
	}
}

func TestGetRoleById(t *testing.T) {
	res, err := GetRoleByUserId(5)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", res)
}

func TestGetRoleByName(t *testing.T) {
	res, err := GetRoleByName("test", "role_1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", res)
}
