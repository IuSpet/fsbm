package pmsg

import (
	"bytes"
	"encoding/json"
	"fsbm/util"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type messagePublicBody struct {
	Method       string `json:"method"`
	AppKey       string `json:"app_key"`
	TargetAppKey string `json:"target_app_key"`
	SignMethod   string `json:"sign_method"`
	Sign         string `json:"sign"`
	Session      string `json:"session"`
	Timestamp    string `json:"timestamp"`
	Format       string `json:"format"`
	V            string `json:"v"`
	PartnerId    string `json:"partner_id"`
	Simplify     bool   `json:"simplify"`
}

type messageReqBody struct {
	Extend          string `json:"extend"`
	SmsType         string `json:"sms_type"`
	SmsFreeSignName string `json:"sms_free_sign_name"`
	SmsParam        string `json:"sms_param"`
	RecNum          string `json:"rec_num"`
	SmsTemplateCode string `json:"sms_template_code"`
}

type messageBody struct {
	messagePublicBody
	messageReqBody
}

type messageRspBody struct {
	AlibabaAliqinFcSmsNumSendResponse struct {
		Result string `json:"result"`
	} `json:"alibaba_aliqin_fc_sms_num_send_response"`
	ErrorResponse struct {
		Msg     string `json:"msg"`
		Code    int64  `json:"code"`
		SubMsg  string `json:"sub_msg"`
		SubCode string `json:"sub_code"`
	} `json:"error_response"`
}

func SendMessage(phone string, alarmParams *util.PhoneMessageModel) error {
	c := &http.Client{}
	param, err := json.Marshal(alarmParams)
	if err != nil {
		return err
	}
	body := newMessageBody(phone, string(param))
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("post", util.PhoneMessageUrl, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("charset", "utf-8")
	rsp, err := c.Do(req)
	defer func() {
		_ = rsp.Body.Close()
	}()
	rspData, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	rspBody := &messageRspBody{}
	err = json.Unmarshal(rspData, rspBody)
	if err != nil {
		return err
	}
	return nil
}

func newMessageBody(param, phone string) *messageBody {
	msgBody := &messageBody{
		messagePublicBody: messagePublicBody{
			Method:       util.PhoneMessageMethod,
			AppKey:       util.PhoneMessageAppKey,
			TargetAppKey: "",
			SignMethod:   "md5",
			Sign:         "",
			Session:      "",
			Timestamp:    time.Now().Format(util.YMDHMS),
			Format:       "json",
			V:            "2",
			PartnerId:    "",
			Simplify:     false,
		},
		messageReqBody: messageReqBody{
			Extend:          "",
			SmsType:         "normal",
			SmsFreeSignName: util.PhoneMessageSignName,
			SmsParam:        param,
			RecNum:          phone,
			SmsTemplateCode: util.PhoneMessageTemplate,
		},
	}
	s := generateSign(msgBody)
	msgBody.Sign = s
	return msgBody
}

type field struct {
	key   string
	value string
}
type fl = []field

func generateSign(body *messageBody) string {
	fieldList := fl{
		{"method", body.Method},
		{"app_key", body.AppKey},
		{"target_app_key", body.AppKey},
		{"sign_method", body.SignMethod},
		{"session", body.Session},
		{"timestamp", body.Timestamp},
		{"format", body.Format},
		{"v", body.V},
		{"partner_id", body.PartnerId},
		{"simplify", strconv.FormatBool(body.Simplify)},
		{"extend", body.Extend},
		{"sms_type", body.SmsType},
		{"sms_free_sign_name", body.SmsFreeSignName},
		{"sms_param", body.SmsParam},
		{"rec_num", body.RecNum},
		{"sms_template_code", body.SmsTemplateCode},
	}
	sort.SliceStable(fieldList, func(i, j int) bool {
		return fieldList[i].key < fieldList[j].key
	})
	s := spliceField(fieldList)
	res := util.Md5(util.PhoneMessageSecret + s + util.PhoneMessageSecret)
	return res
}

func spliceField(f fl) string {
	var res = ""
	for _, item := range f {
		res += item.key
		res += item.value
	}
	return res
}
