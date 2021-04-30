package db

import "time"

var ShopStatusMapping = map[int8]string{
	0: "正常",
	1: "已关闭",
}

type ShopList struct {
	ID           int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	Name         string    `gorm:"varchar(127);not null; comment:店铺名称"`
	UserID       int64     `gorm:"type:bigint;not null; comment:店铺负责人id"`
	Addr         string    `gorm:"type:varchar(255);not null; comment:店铺地址"`
	Latitude     float64   `gorm:"not null; comment:纬度"`
	Longitude    float64   `gorm:"not null; comment:经度"`
	NoticeConfig string    `gorm:"type:varchar(255);not null; comment:店铺报警配置"`
	Status       int8      `gorm:"type:tinyint;not null; comment:状态，0：正常，1：已删除"`
	Remark       string    `gorm:"type:text;;comment:店铺备注"`
	CreatedAt    time.Time `gorm:"autoCreateTime; not null"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime; not null"`
}

type ShopNoticeConfigBase struct {
	IsOn         bool    `json:"is_on"`
	Threshold    int64   `json:"threshold"`
	NoticeDevice []int64 `json:"notice_device"`
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

func SaveShopListRow(row *ShopList) (err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Save(row).Error
	return
}

func GetShopListById(shopIdList []int64) (res []ShopList, err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Where("id in ?", shopIdList).Find(&res).Error
	return
}

func GetShopListByUserId(id int64) (res []ShopList, err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Where("user_id = ?", id).Find(&res).Error
	return
}
