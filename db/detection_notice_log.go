package db

import "time"

// 识别结果报警通知表

type DetectionNoticeLog struct {
	ID        int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	ResultID  int64     `gorm:"type:bigint;comment:识别结果id"`
	UserID    int64     `gorm:"type:bigint;comment:用户Id"`
	Message   string    `gorm:"type:varchar(255);not null;comment:报警内容"`
	NotifyTs  int64     `gorm:"type:bigint;not null; comment: 通知时间"`
	Rank      int8      `gorm:"type:tinyint;not null;comment:报警级别"`
	Status    int8      `gorm:"type:tinyint;not null;comment:日志状态"`
	CreatedAt time.Time `gorm:"autoCreateTime; not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime; not null"`
}

func (DetectionNoticeLog) TableName() string {
	return "shop_list"
}

func init() {
	table := DetectionNoticeLog{}
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
