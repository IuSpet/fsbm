package db

import (
	"fmt"
	"fsbm/conf"
	"testing"
)

func TestSaveUserInfo(t *testing.T) {
	conf.Init()
	Init()
	info := &UserAccountInfo{
		Name:     "luSpet",
		Email:    "1037821259@qq.com",
		Status:   0,
		Password: "123456",
		Phone:    "12345678",
		Gender:   1,
		Age:      22,
	}
	err := SaveUserInfo(info)
	if err != nil {
		panic(err)
	}
}

func TestGetUserByEmail(t *testing.T) {
	conf.Init()
	Init()
	res, err := GetUserByEmail("admin__1@admin.com")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", res)
}

func TestGetAllUser(t *testing.T) {
	conf.Init()
	Init()
	res, err := GetAllUser()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", res)
}

func TestFuzzySearchUser(t *testing.T) {
	conditions := []string{
		"email like %@%",
	}
	res, err := FuzzySearchUser(conditions)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", res)
}

func TestSetAvatar(t *testing.T) {
	conf.Init()
	Init()
	avatar := []byte("123")
	err := SetAvatar("123@321", avatar)
	if err != nil {
		panic(err)
	}
}
