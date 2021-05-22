package cronjob

import (
	"context"
	"fmt"
	"fsbm/db"
	"fsbm/util"
	"fsbm/util/logs"
	"fsbm/util/redis"
	"time"
)

// 监控任务
func recordAlarmCheckTask(ctx context.Context) error {
	fmt.Println("start task alarm")
	// 查询未扫描过的记录
	recordList, err := db.GetUncheckedRecords()
	if err != nil {
		logs.CtxError(ctx, "get unchecked records error. err: %+v", err)
		return err
	}
	recordList = recordsScanNoHat(ctx, recordList)
	err = db.SaveDetectionResultRecordRows(recordList)
	return err
}

// 对未戴帽子扫描
func recordsScanNoHat(ctx context.Context, records []db.DetectionResultRecord) []db.DetectionResultRecord {
	now := time.Now()
	msgTemplate := `
报警店铺：%s，
报警内容：后厨有人员未佩戴帽子，
详细信息："http://47.95.248.242/#/alarm/alarm_detail?id=%d"
请尽快前往查看
`
	alarmCnt := 0
	for idx := range records {
		records[idx].Status = db.DetectionResultRecordStatus_Normal
		if records[idx].NoHatCnt > 0 {
			records[idx].Status = db.DetectionResultRecordStatus_Alarm
			shopList, err := db.GetShopListById([]int64{records[idx].SrcShopID})
			if err != nil || len(shopList) == 0 {
				logs.CtxError(ctx, "get shop info error. err: %+v, shopInfo: %+v", err, shopList)
				continue
			}
			shopInfo := shopList[0]
			// 写报警记录
			alarmRecordRow := &db.RecordAlarm{
				RecordId:  records[idx].ID,
				ShopId:    shopInfo.ID,
				UserId:    shopInfo.UserID,
				MessageId: 0,
				AlarmType: db.RecordAlarmAlarmType_Nohat,
				AlarmAt:   now.Format(util.YMDHMS),
				Status:    0,
			}
			err = db.SaveRecordAlarmRow(alarmRecordRow)
			if err != nil {
				logs.CtxError(ctx, "save record alarm error. err: %+v", err)
				continue
			}
			// 写发送消息
			messageRow := &db.NotifyUserMessage{
				UserId:  shopInfo.UserID,
				Message: fmt.Sprintf(msgTemplate, shopInfo.Name, alarmRecordRow.ID),
				Status:  db.NotifyUserMessageStatus_NotSentYet,
			}
			err = db.SaveNotifyUserMessageRow(messageRow)
			if err != nil {
				logs.CtxError(ctx, "save msg error. err: %+v", err)
				continue
			}
			alarmRecordRow.MessageId = messageRow.ID
			_ = db.SaveRecordAlarmRow(alarmRecordRow)
			alarmCnt += 1
		}
	}
	key := fmt.Sprintf(util.DashboardAlarmCnt, time.Now().Format(util.YMD))
	_ = redis.IncrByWithRetyr(ctx, key, int64(alarmCnt))
	return records
}
