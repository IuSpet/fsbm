package cronjob

import (
	"context"
	"errors"
	"fmt"
	"fsbm/db"
	"fsbm/util"
	"fsbm/util/logs"
	"fsbm/util/mail"
	"fsbm/util/pmsg"
	"fsbm/util/wxmsg"
	"time"
)

// 消息发送任务
func notifyMessageTask(ctx context.Context) error {
	// 未发送消息
	notSentMessageList, err := db.GetNotSentMessageList()
	if err != nil {
		logs.CtxError(ctx, "get not sent msg list error. err: %+v", err)
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
			logs.CtxError(ctx, "[%d]msg's mail send error. err: %+v", message.ID, err)
			if sendThreshold.After(message.CreatedAt) {
				message.Status = db.NotifyUserMessageStatus_AlwaysSentFail
				return errors.New(fmt.Sprintf("[%d]always send fail", message.ID))
			}
		} else {
			message.Status |= db.NotifyUserMessageSendMailSuccess
		}
	}
	// 发送公众号消息
	if message.Status&db.NotifyUserMessageSendWxMessageSuccess == 0 {
		err = wxmsg.SendMsg(user.OpenId, &util.WxMessageModel{
			First:    "食品安全管理后台报警",
			Keyword1: message.Message,
			Keyword2: message.CreatedAt.Format(util.YMDHMS),
			Remark:   "",
		})
		if err != nil {
			logs.CtxError(ctx, "[%d]msg's wx send error. err: %+v", message.ID, err)
			if sendThreshold.After(message.CreatedAt) {
				message.Status = db.NotifyUserMessageStatus_AlwaysSentFail
				return errors.New(fmt.Sprintf("[%d]always send fail", message.ID))
			}
		} else {
			message.Status |= db.NotifyUserMessageSendMailSuccess
		}
	}
	// 发送短信
	if message.Status&db.NotifyUserMessageSendPhoneMessageSuccess == 0 {
		alarm, _ := db.GetAlarmByMessageId(message.ID)
		shop, _ := db.GetShopInfoById(alarm.ShopId)
		err = pmsg.SendMessageV2(user.Phone, &util.PhoneMessageModel{
			ShopName:     shop.Name,
			AlarmContent: db.RecordAlarmAlarmTypeMapping[alarm.AlarmType],
			AlarmDetail:  "请前往系统查看相信信息",
		})
		if err != nil {
			logs.CtxError(ctx, "[%d]msg's wx send error. err: %+v", message.ID, err)
			if sendThreshold.After(message.CreatedAt) {
				message.Status = db.NotifyUserMessageStatus_AlwaysSentFail
				return errors.New(fmt.Sprintf("[%d]always send fail", message.ID))
			}
		} else {
			message.Status |= db.NotifyUserMessageSendPhoneMessageSuccess
		}
	}
	return nil
}
