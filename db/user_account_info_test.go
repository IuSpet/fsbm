package db

import (
	"fmt"
	"testing"
)

func TestSaveUserInfo(t *testing.T) {
	info := &UserAccountInfo{
		Name:     "admin",
		Email:    "admin@admin.com",
		Status:   0,
		Password: "123456",
	}
	err := SaveUserInfo(info)
	if err != nil {
		panic(err)
	}
}

func TestGetUserByEmail(t *testing.T) {
	res, err := GetUserByEmail("admin@admin.com")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", res)
}
