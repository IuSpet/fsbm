package dashboard

import (
	"fmt"
	"fsbm/conf"
	"fsbm/db"
	"testing"
)

func TestGetTodayShopAlarm(t *testing.T) {
	conf.Init()
	db.Init()
	res, err := getTodayShopAlarm()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}

func TestGetShopInfoList(t *testing.T) {
	conf.Init()
	db.Init()
	res, err := getShopInfoList()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}
