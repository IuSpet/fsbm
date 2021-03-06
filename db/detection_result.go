package db

import "time"

type DetectionResult struct {
	ID           int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	SrcVideoID   int64     `gorm:"type:bigint; not null; comment:图片来源视频id"`
	SrcDeviceID  int64     `gorm:"type:bigint; not null; comment:图片来源设别id"`
	At           int64     `gorm:"type:bigint; not null; comment:图片时间戳"`
	FrameCnt     int64     `gorm:"type:bigint; not null; comment:图片在原视频多少帧"`
	Path         string    `gorm:"type:varchar(255); not null; comment:图片存储路径"`
	IdentifyCnt  int64     `gorm:"type:int; not null; comment:总识别结果"`
	WearHatCnt   int64     `gorm:"type:int; not null; comment:戴帽子人数"`
	NoHatCnt     int64     `gorm:"type:int; not null; comment:未戴帽子人数"`
	Log          string    `gorm:"type:varchar(255); not null; comment:报警错误日志"`
	ExternalJson string    `gorm:"type:varchar(255); not null; comment:识别额外信息"`
	CreatedAt    time.Time `gorm:"autoCreateTime; not null"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime; not null"`
}

func (DetectionResult) TableName() string {
	return "detection_result"
}

func init() {
	table := DetectionResult{}
	RegisterMigration(table.TableName(), func() {
		conn, err := fsbmSession.GetConnection()
		if err != nil {
			panic(err)
		}
		err = conn.Debug().Set("gorm:table_options", "ENGINE=INNODB CHARSET=utf8").AutoMigrate(&table)
		if err != nil {
			panic(err)
		}
	})
}

func SaveDetectionResultRows(rows []DetectionResult) (err error) {
	conn, err := fsbmSession.GetConnection()
	if err != nil {
		return err
	}
	err = conn.Debug().Save(rows).Error
	return
}

func GetDetectionResultsByVideoId(id []int64) (res []DetectionResult, err error) {
	conn, err := fsbmSession.GetConnection()
	if err != nil {
		return nil, err
	}
	err = conn.Debug().Where("src_video_id in (?)", id).Find(&res).Error
	return
}
