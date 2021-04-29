package shop

import (
	"fmt"
	"fsbm/conf"
	"fsbm/db"
	"testing"
	"time"
)

func TestGetShopListServer(t *testing.T) {
	conf.Init()
	db.Init()
	begin := time.Now().AddDate(0, 0, -5)
	end := time.Now()
	rows, err := getShopListRows("abc", "123", "xyz", begin, end)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(rows)
}

func TestGetDeviceListServer(t *testing.T) {
	conf.Init()
	db.Init()
	rows, err := getMonitorListRows("abc", "123", "root", "beijing", "hls")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(rows)
}
