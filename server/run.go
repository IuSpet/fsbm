package server

import (
	"fmt"
	"fsbm/db"
	"fsbm/util"
	"github.com/gin-gonic/gin"
)

func Run() {
	InitRoleId()
	router := gin.Default()
	Register(router)
	_ = router.Run()
}

func InitRoleId() {
	admin, err1 := db.GetRoleByName("user", "admin")
	normalUser, err2 := db.GetRoleByName("user", "normal_user")
	supervision, err3 := db.GetRoleByName("user", "supervision")
	shopOwner, err4 := db.GetRoleByName("user", "shop_owner")
	if err1 != nil && err2 != nil && err3 != nil && err4 != nil {
		fmt.Printf("err1: %+v\nerr2: %+v\nerr3: %+v\nerr4: %+v\n", err1, err2, err3, err4)
		panic("init role id error")
	}
	util.Role_AdminId = admin.ID
	util.Role_NormalUserId = normalUser.ID
	util.Role_SupervisionId = supervision.ID
	util.Role_ShopOwnerId = shopOwner.ID
}
