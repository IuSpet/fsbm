package db

import (
	"fsbm/conf"
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
