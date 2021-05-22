package db

import "time"

type UserOperationLog struct {
	ID         int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	UserId     int64     `gorm:"type:bigint;not null;comment:用户id"`
	Operation  string    `gorm:"type:text;not null;comment:操作描述"`
	OperatedAt int64     `gorm:"type:bigint;not null;comment:操作时间"`
	CreatedAt  time.Time `gorm:"autoCreateTime; not null"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime; not null"`
}

func (UserOperationLog) TableName() string {
	return "user_operation_log"
}

func init() {
	table := UserOperationLog{}
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

func SaveUserOperationLog(row *UserOperationLog) (err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Debug().Save(row).Error
	return
}

func GetUserOperationRows(id int64) (res []UserOperationLog, err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Debug().Where("user_id = ?", id).Order("operated_at desc").Find(&res).Error
	return
}
