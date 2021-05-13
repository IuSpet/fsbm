package dashboard

import (
	"fmt"
	"fsbm/db"
	"fsbm/util"
	"fsbm/util/logs"
	"fsbm/util/redis"
	"github.com/gin-gonic/gin"
	"math"
	"time"
)

// 首页数据指标
func GlobalStatsServer(ctx *gin.Context) {
	now := time.Now()
	key1 := fmt.Sprintf(util.DashboardRecordCnt, now.Format(util.YMD))
	recordCnt, err := redis.GetInt64WithRetry(ctx, key1)
	if err != nil {
		logs.CtxError(ctx, "get record cnt error. key: %s, err: %+v", key1, err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	key2 := fmt.Sprintf(util.DashboardAlarmCnt, now.Format(util.YMD))
	alarmCnt, err := redis.GetInt64WithRetry(ctx, key2)
	if err != nil {
		logs.CtxError(ctx, "get alarm cnt error. key: %s, err: %+v", key1, err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	key3 := fmt.Sprintf(util.DashboardLatestRecord, now.Format(util.YMD))
	latestRecord, err := redis.GetWithRetry(ctx, key3)
	if err != nil {
		logs.CtxError(ctx, "get latest record error. key: %s, err: %+v", key1, err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	passRate, err := getShopPassRate()
	if err != nil {
		logs.CtxError(ctx, "get pass rate error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	rsp := dashboardStatsResponse{
		RecordCnt:    recordCnt,
		AlarmCnt:     alarmCnt,
		LatestRecord: latestRecord,
		ShopPassRate: passRate,
	}
	util.EndJson(ctx, rsp)
}

func getShopPassRate() (float64, error) {
	shopList, err := db.GetAvailableShopList()
	if err != nil {
		return 0, err
	}
	alarmShop, err := getTodayShopAlarm()
	if err != nil {
		return 0, err
	}
	totalCnt := len(shopList)
	alarmCnt := len(alarmShop)
	if totalCnt == 0 {
		return 0, nil
	}
	return math.Floor(float64(alarmCnt)/float64(totalCnt)*10000.0) / 10000.0, nil
}
