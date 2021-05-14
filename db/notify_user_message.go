package db

import "time"

const (
	NotifyUserMessageSendMailSuccess         int8 = 1
	NotifyUserMessageSendPhoneMessageSuccess int8 = 2
	NotifyUserMessageSendWxMessageSuccess    int8 = 4
	NotifyUserMessageStatus_NotSentYet       int8 = 0
	NotifyUserMessageStatus_HasSent               = NotifyUserMessageSendMailSuccess | NotifyUserMessageSendPhoneMessageSuccess | NotifyUserMessageSendWxMessageSuccess
	NotifyUserMessageStatus_AlwaysSentFail   int8 = -1
)

var NotifyUserMessageStatusMapping = map[int8]string{
	NotifyUserMessageStatus_NotSentYet:     "未发送",
	NotifyUserMessageStatus_HasSent:        "发送成功",
	NotifyUserMessageStatus_AlwaysSentFail: "发送失败(不再重试)",
}

type NotifyUserMessage struct {
	ID        int64     `gorm:"AUTO_INCREMENT; primaryKey"`
	UserId    int64     `gorm:"type:bigint;not null;comment:用户id"`
	Message   string    `gorm:"type:text;not null;comment:消息内容"`
	Status    int8      `gorm:"type:tinyint;not null;comment:消息发送状态"`
	SentAt    int64     `gorm:"type:bigint;消息成功发送时间戳"`
	CreatedAt time.Time `gorm:"autoCreateTime; not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime; not null"`
}

func (NotifyUserMessage) TableName() string {
	return "notify_user_message"
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

func SaveNotifyUserMessageRow(row *NotifyUserMessage) (err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Save(row).Error
	return
}

func SaveNotifyUserMessageRows(rows []NotifyUserMessage) (err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	err = conn.Save(rows).Error
	return
}

func GetNotSentMessageList() (res []NotifyUserMessage, err error) {
	conn, err := FsbmSession.GetConnection()
	if err != nil {
		return
	}
	statusList := []int8{
		NotifyUserMessageStatus_NotSentYet,
		NotifyUserMessageSendMailSuccess,
		NotifyUserMessageSendPhoneMessageSuccess,
		NotifyUserMessageSendWxMessageSuccess,
		NotifyUserMessageSendMailSuccess | NotifyUserMessageSendPhoneMessageSuccess,
		NotifyUserMessageSendMailSuccess | NotifyUserMessageSendWxMessageSuccess,
		NotifyUserMessageSendPhoneMessageSuccess | NotifyUserMessageSendWxMessageSuccess,
	}
	err = conn.Where("status in ?", statusList).Find(res).Error
	return
}
