package cronjob

import "time"

var CronTaskConfig = map[string]TaskExecutor{
	"auth_check": NewTickerTask("auth_check", 2*time.Hour, authCheckTask),
}

func RunCronJob() {
	for _, task := range CronTaskConfig {
		go task.Execute()
	}
}
