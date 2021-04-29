package db

import (
	"fsbm/conf"
	"testing"
)

func TestSaveMonitorListRow(t *testing.T) {
	conf.Init()
	Init()
	row := &MonitorList{
		ShopId:    0,
		Name:      "test_monitor_1",
		VideoType: "hls",
		VideoSrc:  "//sf1-cdn-tos.huoshanstatic.com/obj/media-fe/xgplayer_doc_video/hls/xgplayer-demo.m3u8",
	}
	err := SaveMonitorListRow(row)
	if err != nil {
		t.Error(err)
	}
}
