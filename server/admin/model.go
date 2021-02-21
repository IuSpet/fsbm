package admin

type getUserListRequest struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
}

type getUserListResponse struct {
	UserInfoList []userInfo `json:"user_info_list"`
	TotalCount   int64      `json:"total_count"`
}

type userInfo struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type getUserDetailRequest struct {
	Email string `json:"email"`
}

type getUserDetailResponse struct {
	Email  string           `json:"email"`
	Name   string           `json:"name"`
	Status string           `json:"status"`
	Roles  []userDetailRole `json:"roles"`
}

type userDetailRole struct {
	Type string `json:"type"`
	Name string `json:"name"`
}
