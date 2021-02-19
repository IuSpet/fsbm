package db

import (
	"fmt"
	"testing"
)

func TestAddRole(t *testing.T) {
	role := &AuthRole{
		Role:   "role_2",
		Type:   "test",
		Status: 1,
	}
	err := SaveRole(role)
	if err != nil {
		panic(err)
	}
}

func TestGetRoleById(t *testing.T) {
	res, err := GetRoleById(5)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", res)
}
