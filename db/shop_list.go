package db

import "time"

type ShopList struct {
	ID         int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	Admin      string    `gorm:"type:varchar(127);not null; comment:店铺负责人"`
	AdminEmail string    `gorm:"type:varchar(127);not null; comment:店铺负责人邮箱"`
	AdminPhone string    `gorm:"type:varchar(127);not null; comment:店铺负责人手机"`
	Addr       string    `gorm:"type:varchar(255);not null; comment:店铺地址"`
	CreatedAt  time.Time `gorm:"autoCreateTime; not null"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime; not null"`
}

func (ShopList) TableName() string {
	return "shop_list"
}

func init() {
	table := ShopList{}
	RegisterMigration(table.TableName(), func() {
		conn, err := FsbmSession.GetConnection()
		if err != nil {
			panic(err)
		}
		err = conn.Set("gorm:table_options", "ENGINE=INNODB CHARSET=utf8").AutoMigrate(&table)
		if err != nil {
			panic(err)
		}
	})
}
