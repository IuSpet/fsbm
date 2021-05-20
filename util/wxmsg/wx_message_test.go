package wxmsg

import (
	"fsbm/util"
	"testing"
	"time"
)

func TestSendMsg(t *testing.T) {
	err := SendMsg("oUip_t8R1AzvXQDjDb5AzZIX4by4", &util.WxMessageModel{
		First:    "测试报警标题",
		Keyword1: "未戴帽子",
		Keyword2: time.Now().Format(util.YMDHMS),
		Remark:   "测试",
	})
	if err != nil {
		panic(err)
	}

}
