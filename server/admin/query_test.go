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
	rows, err := getUserList("abc", "", "", 1, -1, time.Now(), time.Now().AddDate(0, 0, 1), 10, 20)
	if err != nil {
		panic(err)
	}
	fmt.Println(rows)
}
