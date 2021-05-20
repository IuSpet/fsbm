package cronjob

import (
	"fmt"
	"time"
)

var CronTaskConfig = map[string]*TaskExecutor{
	"auth_check":         NewTickerTask("auth_check", 2*time.Hour, authCheckTask),
	"record_alarm_check": NewTickerTask("record_alarm_check", 5*time.Minute, recordAlarmCheckTask),
	"notify_message":     NewTickerTask("notify_message", 2*time.Minute, notifyMessageTask),
}

func RunCronJob() {
	for _, task := range CronTaskConfig {
		//fmt.Println(task)
		fmt.Println(task.GetName())
		go task.Execute()
	}
}
