package db

import "time"

// 店铺报警记录表
type ShopNoticeRecord struct {
	ID        int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	UserID    int64     `gorm:"type:bigint;not null; comment:用户ID"`
	ShopID    int64     `gorm:"type:bigint;not null; comment:店铺ID"`
	TaskAt    int64     `gorm:"type:bigint;not null; comment:任务时间"`
	NoticeAt  int64     `gorm:"type:bigint;not null; comment: 报警时间"`
	CreatedAt time.Time `gorm:"autoCreateTime; not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime; not null"`
}

func (ShopNoticeRecord) TableName() string {
	return "shop_notice_record"
}
