package dashboard

import (
	"context"
	"fmt"
	"fsbm/conf"
	"fsbm/db"
	"fsbm/util"
	"fsbm/util/redis"
	"testing"
	"time"
)

func TestGetRedisValue(t *testing.T) {
	now := time.Now()
	ctx := context.Background()
	conf.Init()
	db.Init()
	key1 := fmt.Sprintf(util.DashboardRecordCnt, now.Format(util.YMD))
	recordCnt, err := redis.GetInt64WithRetry(ctx, key1)
	if err != nil {
		panic(err)
	}
	fmt.Println(recordCnt)
	key2 := fmt.Sprintf(util.DashboardAlarmCnt, now.Format(util.YMD))
	alarmCnt, err := redis.GetInt64WithRetry(ctx, key2)
	if err != nil {
		panic(err)
	}
	fmt.Println(alarmCnt)
	key3 := fmt.Sprintf(util.DashboardLatestRecord, now.Format(util.YMD))
	latestRecord, err := redis.GetWithRetry(ctx, key3)
	if err != nil {
		panic(err)
	}
	fmt.Println(latestRecord)
}

func TestSetValue(t *testing.T) {
	now := time.Now()
	ctx := context.Background()
	conf.Init()
	db.Init()
	key1 := fmt.Sprintf(util.DashboardRecordCnt, now.Format(util.YMD))
	redis.SetWithRetry(ctx, key1, 8, 0)
	key2 := fmt.Sprintf(util.DashboardAlarmCnt, now.Format(util.YMD))
	redis.SetWithRetry(ctx, key2, 2, 0)
}

func TestShopPassRate(t *testing.T) {
	conf.Init()
	db.Init()
	res, err := getShopPassRate()
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
