package cronjob

import (
	"context"
	"fmt"
	"fsbm/util"
	"testing"
	"time"
)

func TestNewTickerTask(t *testing.T) {
	task := NewTickerTask("123", 10*time.Second, func(ctx context.Context) error {
		fmt.Println(time.Now().Format(util.YMDHMS))
		return nil
	})
	go task.Execute()
	time.Sleep(30*time.Second)
}
