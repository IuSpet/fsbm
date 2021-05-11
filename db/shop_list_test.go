package db

import (
	"fsbm/conf"
	"math/rand"
	"strconv"
	"testing"
)

func TestSaveShopListRow(t *testing.T) {
	conf.Init()
	Init()
	row := &ShopList{
		Name:         "test_shop1",
		UserID:       8,
		Addr:         "beijing",
		NoticeConfig: "{}",
		Status:       0,
		Remark:       "test shop 1",
	}
	err := SaveShopListRow(row)
	if err != nil {
		t.Error(err)
	}
}

func TestInsertTestRows(t *testing.T) {
	conf.Init()
	Init()
	for i := 0; i < 25; i++ {
		row := &ShopList{
			Name:         "test_shop_" + strconv.FormatInt(int64(i), 10),
			UserID:       int64(rand.Intn(20)),
			Addr:         "test",
			Latitude:     float64(rand.Intn(1800) / 100.0),
			Longitude:    float64(rand.Intn(1800) / 100.0),
			NoticeConfig: "{}",
			Status:       0,
			Remark:       "",
		}
		SaveShopListRow(row)
	}
}
