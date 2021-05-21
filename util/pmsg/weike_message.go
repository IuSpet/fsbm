package pmsg

import (
	"context"
	"encoding/json"
	"fmt"
	"fsbm/util"
	"fsbm/util/logs"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func SendMessageV2(phone string, param *util.PhoneMessageModel) error {
	c := &http.Client{}
	ctx := context.Background()
	b, err := generateParam(phone, param)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, util.PhoneMessageV2Url, b)
	if err != nil {
		logs.CtxError(ctx, "%+v", err)
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("charset", "utf-8")
	req.Header.Add("x-requested-with", "remote-post")
	rsp, err := c.Do(req)
	if err != nil {
		logs.CtxError(ctx, "%+v", err)
		return err
	}
	defer func() {
		_ = rsp.Body.Close()
	}()
	rspData, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		logs.CtxError(ctx, "%+v", err)
		return err
	}
	logs.CtxInfo(ctx, "rsp: %s", string(rspData))
	return nil
}

func generateParam(phone string, param *util.PhoneMessageModel) (io.Reader, error) {
	s, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(s))
	params := url.Values{}
	params["action"] = []string{util.PhoneMessageV2Action}
	params["template_code"] = []string{util.PhoneMessageV2TemplateCode}
	params["template_param"] = []string{string(s)}
	params["mobile"] = []string{phone}
	fmt.Println(params.Encode())
	return strings.NewReader(params.Encode()), nil
}
