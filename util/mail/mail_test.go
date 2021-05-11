package mail

import "testing"

func TestSendAlarmMail(t *testing.T) {
	content := `
【食品安全管理后台报警】
报警店铺：水天堂工农路店
报警内容：后厨有人员未佩戴帽子
详细信息："http://localhost:9528/#/alarm/alarm_detail?id=3"
请尽快前往查看
`
	msg := &DefaultMail{
		Dest:    []string{"1037821259@qq.com"},
		Subject: "食品安全管理后台报警",
		Text:    []byte(content),
	}
	err := SendMail(msg)
	if err != nil {
		t.Error(err)
	}
}
