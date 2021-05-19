package detection

import (
	"fsbm/conf"
	"fsbm/db"
	"testing"
)

func TestGetShopInfo(t *testing.T){
	conf.Init()
	db.Init()
	_,_ = getShopInfo("aaa","bbb","ccc")
}