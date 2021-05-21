package alarm

import (
	"fmt"
	"fsbm/conf"
	"fsbm/db"
	"testing"
)

func TestAlarmDetailInfo(t *testing.T) {
	conf.Init()
	db.Init()
	res, err := getAlarmInfo(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", res)
}
