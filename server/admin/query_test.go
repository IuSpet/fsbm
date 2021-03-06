package admin

import (
	"fmt"
	"fsbm/conf"
	"fsbm/db"
	"testing"
	"time"
)

func TestGetUserList(t *testing.T) {
	conf.Init()
	db.Init()
	rows, err := getUserList("", "", "", -1, -1, time.Unix(0, 0), time.Now().AddDate(0, 0, 1))
	if err != nil {
		panic(err)
	}
	fmt.Println(rows)
}
