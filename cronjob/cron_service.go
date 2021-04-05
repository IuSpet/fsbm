package cronjob

import "context"

type MonitorTask interface {
	Run()
}

// 用cron表达式定义的计划任务
type CronTask struct {
	Name string
	Ctx  context.Context
	Task func()
}

// 定时间隔执行的计划任务
type TickerTask struct {
}
