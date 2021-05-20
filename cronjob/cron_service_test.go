package cronjob

import (
	"context"
	"fmt"
	"fsbm/util"
	"testing"
	"time"
)

func TestNewTickerTask(t *testing.T) {
	task1 := NewTickerTask("123", 10*time.Second, func(ctx context.Context) error {
		fmt.Println(time.Now().Format(util.YMDHMS))
		return nil
	})
	task2 := NewTickerTask("321",7 * time.Second, func(ctx context.Context) error {
		fmt.Println("this is task2")
		return nil
	})
	var tasks = map[string]*TaskExecutor{
		"t1":task1,
		"t2":task2,
	}
	for _,task := range tasks{
		go task.Execute()
	}
	time.Sleep(30*time.Second)
}
