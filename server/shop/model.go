package shop

import (
	"fsbm/db"
	"fsbm/util"
	"time"
)

type getShopListRequest struct {
	Name        string           `json:"name"`
	Admin       string           `json:"admin"`
	Addr        string           `json:"addr"`
	CreateBegin string           `json:"create_begin"`
	CreateEnd   string           `json:"create_end"`
	SortFields  []util.SortField `json:"sort_fields"`
	Page        int64            `json:"page"`
	PageSize    int64            `json:"page_size"`
}

type shopInfo struct {
	Name       string `json:"name"`
	AdminName  string `json:"admin_name"`
	AdminPhone string `json:"admin_phone"`
	AdminEmail string `json:"admin_email"`
	Addr       string `json:"addr"`
	CreatedAt  string `json:"created_at"`
	Status     string `json:"status"`
}

type getShopListResponse struct {
	List     []shopInfo `json:"list"`
	TotalCnt int64      `json:"total_cnt"`
}

type shopInfoRow struct {
	Name       string
	AdminName  string
	AdminPhone string
	AdminEmail string
	Addr       string
	CreatedAt  time.Time
	Status     int8
}

type addShopRequest struct {
	Name      string                             `json:"name"`
	Addr      string                             `json:"addr"`
	Latitude  float64                            `json:"latitude"`
	Longitude float64                            `json:"longitude"`
	NoticeCfg map[string]db.ShopNoticeConfigBase `json:"notice_cfg"`
	Remark    string                             `json:"remark"`
}

type getDeviceListRequest struct {
	DeviceName string `json:"device_name"`
	ShopName   string `json:"shop_name"`
	AdminName  string `json:"admin_name"`
	Addr       string `json:"addr"`
	VideoType  string `json:"video_type"`
	util.ListReqField
}

type getMonitorListResponse struct {
	List     []monitorInfo `json:"list"`
	TotalCnt int64         `json:"total_cnt"`
}

type monitorInfo struct {
	MonitorName string `gorm:"column:monitor_name" json:"monitor_name"`
	ShopName    string `gorm:"column:shop_name" json:"shop_name"`
	Addr        string `gorm:"column:addr" json:"addr"`
	AdminName   string `gorm:"column:user_name" json:"admin_name"`
	AdminPhone  string `gorm:"column:user_phone" json:"admin_phone"`
	VideoType   string `gorm:"column:video_type" json:"video_type"`
	VideoSrc    string `gorm:"column:video_src" json:"video_src"`
}

type getLiveWallSrcResponse struct {
	List []liveSrcInfo `json:"list"`
}

type liveSrcInfo struct {
	MonitorName string `json:"monitor_name"`
	ShopName    string `json:"shop_name"`
	VideoType   string `json:"video_type"`
	VideoSrc    string `json:"video_src"`
}

type addMonitorRequest struct {
	ShopId    int64  `json:"shop_id"`
	Name      string `json:"name"`
	VideoType string `json:"video_type"`
	VideoSrc  string `json:"video_src"`
}

type getShopByEmailRequest struct {
	UserEmail string `json:"user_email"`
}

type getShopByEmailResponse struct {
	List []userShopInfo `json:"list"`
}

type userShopInfo struct {
	ShopId   int64  `json:"shop_id"`
	ShopName string `json:"shop_name"`
	Addr     string `json:"addr"`
}
