package task

import (
	"errors"
	"fsbm/db"
)

func GetUserIdByShopId(shopId int64) (int64, error) {
	shopInfo, err := db.GetShopListById([]int64{shopId})
	if err != nil {
		return -1, err
	}
	if len(shopInfo) < 0 {
		return -1, errors.New("no shop result")
	}
	return shopInfo[0].UserID, nil
}
