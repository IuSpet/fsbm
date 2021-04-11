package admin

import (
	"fmt"
	"fsbm/conf"
	"fsbm/db"
	"testing"
)

func TestGetSortedUserList(t *testing.T) {
	conf.Init()
	db.Init()
	req := newGetUserListRequest()
	req.SortFields = []sortField{
		{
			Field: "Name",
			Order: "asc",
		},
	}
	res, cnt, err := getSortedUserList(&req, false)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", res)
	fmt.Println(cnt)
}
