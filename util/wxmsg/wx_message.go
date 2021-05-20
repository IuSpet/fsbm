package wxmsg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fsbm/util"
	"io/ioutil"
	"net/http"
)

type dataValue struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

type messageBody struct {
	Touser     string `json:"touser"`
	TemplateId string `json:"template_id"`
	Url        string `json:"url"`
	Data       struct {
		First    dataValue `json:"first"`
		Keyword1 dataValue `json:"keyword1"`
		Keyword2 dataValue `json:"keyword2"`
		Remark   dataValue `json:"remark"`
	} `json:"data"`
}

func SendMsg(openId string, msg *util.WxMessageModel) error {
	c := &http.Client{}
	body := newMessageBody(openId, msg)
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("post", util.WxMessageUrl, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("charset", "utf-8")
	rsp, err := c.Do(req)
	defer func() {
		_ = rsp.Body.Close()
	}()
	if err != nil {
		return err
	}
	rspData, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(rspData))
	return nil
}

func newMessageBody(id string, param *util.WxMessageModel) *messageBody {
	return &messageBody{
		Touser:     id,
		TemplateId: util.WxMessageTemplate,
		Url:        "",
		Data: struct {
			First    dataValue `json:"first"`
			Keyword1 dataValue `json:"keyword1"`
			Keyword2 dataValue `json:"keyword2"`
			Remark   dataValue `json:"remark"`
		}{
			First: dataValue{
				Value: param.First,
				Color: "",
			},
			Keyword1: dataValue{
				Value: param.Keyword1,
				Color: "",
			},
			Keyword2: dataValue{
				Value: param.Keyword2,
				Color: "",
			},
			Remark: dataValue{
				Value: param.Remark,
				Color: "",
			},
		},
	}
}
