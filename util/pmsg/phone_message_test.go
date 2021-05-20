package pmsg

import (
	"fsbm/util"
	"testing"
)

func TestGenerateSign(t *testing.T) {
	generateSign(&messageBody{})
}

func TestSendMessage(t *testing.T) {
	model := &util.PhoneMessageModel{
		ShopName:     "test_shop",
		AlarmContent: "no hat",
		AlarmDetail:  "url",
	}
	err := SendMessage("18512807827", model)
	if err != nil {
		t.Error(err)
	}
}
