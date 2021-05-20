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
	ShopName     string `json:"shop_name"`
	AlarmContent string `json:"alarm_content"`
	AlarmDetail  string `json:"alarm_detail"`
}
