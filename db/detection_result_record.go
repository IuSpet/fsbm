package db

import "time"

const (
	DetectionResultRecordStatus_NotScanYet int8 = 0 // 未扫描
	DetectionResultRecordStatus_Normal     int8 = 1 // 扫描正常
	DetectionResultRecordStatus_Alarm      int8 = 2 // 扫描触发报警条件
)

type DetectionResultRecord struct {
	ID          int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	SrcDeviceID int64     `gorm:"type:bigint; not null; comment:图片来源设别id"`
	SrcShopID   int64     `gorm:"type:bigint; not null; comment:来源店铺id;"`
	VideoPath   string    `gorm:"type:varchar(255); not null; comment:视频存储路径"`
	At          int64     `gorm:"type:bigint; not null; comment:图片时间戳"`
	FrameCnt    int64     `gorm:"type:bigint; not null; comment:图片在原视频多少帧"`
	ImgPath     string    `gorm:"type:varchar(255); not null; comment:图片存储路径"`
	IdentifyCnt int64     `gorm:"type:int; not null; comment:总识别结果"`
	WearHatCnt  int64     `gorm:"type:int; not null; comment:戴帽子人数"`
	NoHatCnt    int64     `gorm:"type:int; not null; comment:未戴帽子人数"`
	ExtraJson   string    `gorm:"type:varchar(255); not null; comment:识别额外信息"`
	Status      int8      `gorm:"type:tinyint;not null; comment:状态"`
	CreatedAt   time.Time `gorm:"autoCreateTime; not null"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime; not null"`
}
func (DetectionResultRecord) TableName() string {
	return "detection_result"
}

func init() {
	table := DetectionResultRecord{}
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