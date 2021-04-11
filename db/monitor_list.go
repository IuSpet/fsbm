package db

import "time"

type MonitorList struct {
	ID        int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	ShopId    int64     `gorm:"type:bigint;not null; comment: 所属店铺id"`
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
