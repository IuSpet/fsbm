package detection

import "fsbm/db"

func getShopInfo(email, userName, shopName string) (res []db.ShopList, err error) {
	conn, err := db.FsbmSession.GetConnection()
	if err != nil {
		return
	}
	conn = conn.Table("shop_list a " +
		"left join user_account_info b on a.user_id = b.id")
	if email != "" {
		conn = conn.Where("b.email = ?", email)
	}
	if userName != "" {
		conn = conn.Where("b.name = ?", userName)
	}
	if shopName != "" {
		conn = conn.Where("a.name = ?", shopName)
	}
	err = conn.Debug().Find(&res).Error
	return
}
