package cronjob

import (
	"context"
	"fsbm/util/logs"
	"time"
)



type TaskExecutor struct {
	context.Context
	MonitorTask
	sig <-chan int
}

type MonitorTask interface {
	SetTrigger() <-chan int
	GetName() string
	Run(context.Context) error
}

// 用cron表达式定义的计划任务
type CronTask struct {
	Name string
	Task func(context.Context) error
	cron string
}

// 定时间隔执行的计划任务
type TickerTask struct {
	Name     string
	Task     func(context.Context) error
	interval time.Duration
}

func (e *TaskExecutor) Execute() {
	for {
		select {
		case <-e.sig:
			logs.CtxInfo(e, "start task: %s", e.GetName())
			err := e.Run(e)
			if err != nil {
				logs.CtxError(e, "task execute fail. err: %+v", err)
			}
		}
	}
}

func (t *TickerTask) SetTrigger() <-chan int {
	ch := make(chan int, 1)
	go func() {
		for {
			ch <- 1
			time.Sleep(t.interval)
		}
	}()
	return ch
}

func (t *TickerTask) GetName() string {
	return t.Name
}

func (t *TickerTask) Run(ctx context.Context) error {
	return t.Task(ctx)
}

func NewTickerTask(name string, d time.Duration, fn func(context.Context) error) TaskExecutor {
	tickerTask := TickerTask{
		Name:     name,
		Task:     fn,
		interval: d,
	}
	return TaskExecutor{
		Context:     context.WithValue(context.Background(), "task_name", tickerTask.Name),
		MonitorTask: &tickerTask,
		sig:         tickerTask.SetTrigger(),
	}
}