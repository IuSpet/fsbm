package cronjob

import (
	"context"
	"errors"
	"fmt"
	"fsbm/db"
	"fsbm/util/logs"
	"fsbm/util/mail"
	"time"
)

// 消息发送任务
func notifyMessageTask(ctx context.Context) error {
	// 未发送消息
	notSentMessageList, err := db.GetNotSentMessageList()
	if err != nil {
		logs.CtxError(ctx, "get not sent message list error. err: %+v", err)
		return err
	}
	for idx := range notSentMessageList {
		if err = sendMessageAllChannel(ctx, &notSentMessageList[idx]); err != nil {
			logs.CtxError(ctx, "err: %+v", err)
		}
	}
	err = db.SaveNotifyUserMessageRows(notSentMessageList)
	return err
}

func sendMessageAllChannel(ctx context.Context, message *db.NotifyUserMessage) error {
	sendThreshold := time.Now().AddDate(0, 0, -1)
	user, err := db.GetUserAccountInfoById(message.UserId)
	if err != nil || user == nil {
		logs.CtxError(ctx, "err: %+v, user: %+v", err, user)
		return err
	}
	// 发送邮件
	if message.Status&db.NotifyUserMessageSendMailSuccess == 0 {
		err = mail.SendMail(&mail.DefaultMail{
			Dest:    []string{user.Email},
			Subject: "【食品安全管理后台报警】",
			Text:    []byte(message.Message),
		})
		if err != nil {
			logs.CtxError(ctx, "[%d]message's mail send error. err: %+v", message.ID, err)
			if sendThreshold.After(message.CreatedAt) {
				message.Status = db.NotifyUserMessageStatus_AlwaysSentFail
				return errors.New(fmt.Sprintf("[%d]always send fail", message.ID))
			}
		} else {
			message.Status |= db.NotifyUserMessageSendMailSuccess
		}
	}
	// 发送短信
	return nil
}
