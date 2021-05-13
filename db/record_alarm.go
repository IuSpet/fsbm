package db

import "time"

const (
	RecordAlarmAlarmType_Nohat int8 = 1 // 未戴厨师帽
)

var RecordAlarmAlarmTypeMapping = map[int8]string{
	RecordAlarmAlarmType_Nohat: "未戴厨师帽",
}

// 报警记录表
type RecordAlarm struct {
	ID        int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	RecordId  int64     `gorm:"type:bigint;not null; comment:记录id"`
	ShopId    int64     `gorm:"type:bigint;not null; comment:店铺id"`
	UserId    int64     `gorm:"type:bigint;not null; comment:用户id"`
	MessageId int64     `gorm:"type:bigint;not null; comment:消息id"`
	AlarmType int8      `gorm:"type:tinyint;not null; comment:报警类型（报警原因）"`
	AlarmAt   string    `gorm:"type:varchar(63);not null; comment:报警时间"`
	Status    int8      `gorm:"type:tinyint;not null; comment:报警状态"`
	CreatedAt time.Time `gorm:"autoCreateTime; not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime; not null"`
}

func (RecordAlarm) TableName() string {
	return "record_alarm"
}

func init() {
	table := RecordAlarm{}
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

func SaveRecordAlarmRow(row *RecordAlarm) (err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Save(row).Error
	return
}

func SaveRecordAlarmRows(rows []RecordAlarm) (err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Save(rows).Error
	return
}
