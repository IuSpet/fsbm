package util

type ListReqField struct {
	Page       int64       `json:"page"`
	PageSize   int64       `json:"page_size"`
	SortFields []SortField `json:"sort_fields"`
}

type SortField struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

type PhoneMessageModel struct {
	ShopName     string `json:"shopName"`
	AlarmContent string `json:"alarmContent"`
	AlarmDetail  string `json:"alarmDetail"`
}

type WxMessageModel struct {
	First    string `json:"first"`
	Keyword1 string `json:"keyword1"`
	Keyword2 string `json:"keyword2"`
	Remark   string `json:"remark"`
}
