package db

import (
	"fsbm/util"
	"time"
)

const (
	_INVALID = -1000 // 无效记录
	_NEW     = 0     // 新记录
	_INFO    = 1     // 普通报警扫描过
	_WARN    = 2     // 警告报警扫描过
	_ERROR   = 4     // 错误报警扫描过
)

const totalStatus = 1 << 3

type DetectionResult struct {
	ID          int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	SrcVideoID  int64     `gorm:"type:bigint; not null; comment:图片来源视频id"`
	SrcDeviceID int64     `gorm:"type:bigint; not null; comment:图片来源设别id"`
	SrcShopID   int64     `gorm:"type:bigint not null; comment:来源店铺id;"`
	At          int64     `gorm:"type:bigint; not null; comment:图片时间戳"`
	FrameCnt    int64     `gorm:"type:bigint; not null; comment:图片在原视频多少帧"`
	Path        string    `gorm:"type:varchar(255); not null; comment:图片存储路径"`
	IdentifyCnt int64     `gorm:"type:int; not null; comment:总识别结果"`
	WearHatCnt  int64     `gorm:"type:int; not null; comment:戴帽子人数"`
	NoHatCnt    int64     `gorm:"type:int; not null; comment:未戴帽子人数"`
	Log         string    `gorm:"type:varchar(255); not null; comment:报警错误日志"`
	ExtraJson   string    `gorm:"type:varchar(255); not null; comment:识别额外信息"`
	Status      int8      `gorm:"type:tinyint;not null; comment:状态"`
	CreatedAt   time.Time `gorm:"autoCreateTime; not null"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime; not null"`
}

func (DetectionResult) TableName() string {
	return "detection_result"
}

func init() {
	table := DetectionResult{}
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

func SaveDetectionResultRows(rows []DetectionResult) (err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return err
	}
	err = conn.Debug().Save(rows).Error
	return
}

func GetDetectionResultsByVideoId(id []int64) (res []DetectionResult, err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return nil, err
	}
	err = conn.Debug().Where("src_video_id in (?)", id).Find(&res).Error
	return
}

func GetDetectionResultsByNoticeLevel(lv util.NoticeLevel) (res []DetectionResult, err error) {
	var statusRange []int
	statusRange = append(statusRange, 0)
	switch lv {
	case util.InfoNotice:
		for i := 0; i < totalStatus; i++ {
			if i&_INFO > 0 {
				statusRange = append(statusRange, i)
			}
		}
	case util.WarnNotice:
		for i := 0; i < totalStatus; i++ {
			if i&_WARN > 0 {
				statusRange = append(statusRange, i)
			}
		}
	case util.ErrorNotice:
		for i := 0; i < totalStatus; i++ {
			if i&_ERROR > 0 {
				statusRange = append(statusRange, i)
			}
		}
	}
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Debug().Where("status in ?", statusRange).Find(&res).Error
	return
}
