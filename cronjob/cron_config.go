package cronjob

import "time"

var CronTaskConfig = map[string]TaskExecutor{
	"auth_check":         NewTickerTask("auth_check", 2*time.Hour, authCheckTask),
	"record_alarm_check": NewTickerTask("record_alarm_check", 5*time.Minute, recordAlarmCheckTask),
	"notify_message":     NewTickerTask("notify_message", 2*time.Minute, notifyMessageTask),
}

func RunCronJob() {
	for _, task := range CronTaskConfig {
		go task.Execute()
	}
}
