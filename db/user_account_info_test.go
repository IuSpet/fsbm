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
		Name:     "test_save3",
		Email:    "test_save3@qq.com",
		Status:   0,
		Password: "test123456",
		Phone:    "12345678",
		Gender:   1,
		Age:      22,
	}
	err := SaveUserInfo(info)
	if err != nil {
		panic(err)
	}
	fmt.Println(info)
	info.Name = "modify_test_save3"
	err = SaveUserInfo(info)
	if err != nil {
		panic(err)
	}
	fmt.Println(info)
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

func TestGetUserAccountInfoById(t *testing.T) {
	conf.Init()
	Init()
	res, err := GetUserAccountInfoById(19)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
