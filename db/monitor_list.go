package db

import "time"

type MonitorList struct {
	ID        int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	ShopId    int64     `gorm:"type:bigint;not null; comment: 所属店铺id"`
	Name      string    `gorm:"type:varchar(127); not null; comment:监控名"`
	VideoType string    `gorm:"type:varchar(63); not null; comment:flv、hls"`
	VideoSrc  string    `gorm:"type:varchar(255); not null; comment:视频源"`
	CreatedAt time.Time `gorm:"autoCreateTime; not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime; not null"`
}

func (MonitorList) TableName() string {
	return "monitor_list"
}

func init() {
	table := MonitorList{}
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

func SaveMonitorListRow(row *MonitorList) (err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Save(row).Error
	return
}

func GetLiveMonitorRows() (res []MonitorList, err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Where("video_src <> ''").Find(&res).Error
	return
}

func GetMonitorListByShopId(id int64) (res []MonitorList, err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Where("shop_id = ?", id).Find(&res).Error
	return
}
