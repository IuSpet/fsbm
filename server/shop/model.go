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
	Latitude  int                                `json:"latitude"`
	Longitude int                                `json:"longtitude"`
	NoticeCfg map[string]db.ShopNoticeConfigBase `json:"notice_cfg"`
	Remark    string                             `json:"remark"`
}
