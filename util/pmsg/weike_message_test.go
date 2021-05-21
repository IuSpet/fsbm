package pmsg

import (
	"fmt"
	"fsbm/util"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestSendMessageV2(t *testing.T) {
	err := SendMessageV2("18512807827", &util.PhoneMessageModel{
		ShopName:     "小小食其家2",
		AlarmContent: "有老鼠进来了2",
		AlarmDetail:  "视频识别发现老鼠2",
	})
	if err != nil {
		panic(err)
	}
}

func TestPostManReq(t *testing.T){
	url := "http://www.weikebaijia.net/ylxxt/sms_message_servlet_action"
	method := "POST"

	payload := strings.NewReader("action=send_template_sms_message&template_code=SMS_217415565&template_param=%7B%22shopName%22%3A%22%E5%B0%8F%E5%B0%8F%E9%A3%9F%E5%85%B6%E5%AE%B6%22%2C%22alarmContent%22%3A%22%E6%9C%89%E8%80%81%E9%BC%A0%E8%BF%9B%E6%9D%A5%E4%BA%86%EF%BC%81%22%2C%22alarmDetail%22%3A%22%E8%A7%86%E9%A2%91%E8%AF%86%E5%88%AB%E5%8F%91%E7%8E%B0%E8%80%81%E9%BC%A0%22%7D&mobile=18512807827")

	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-requested-with", "remote-post")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "JSESSIONID=BA80C23F2DE3416D17790E5B92CA0D2A")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}