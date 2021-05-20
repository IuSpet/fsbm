package admin

type getUserListRequest struct {
	Name        string      `json:"name"`
	Gender      int8        `json:"gender"`
	Age         int8        `json:"age"`
	CreateBegin string      `json:"create_begin"`
	CreateEnd   string      `json:"create_end"`
	Email       string      `json:"email"`
	Phone       string      `json:"phone"`
	Page        int64       `json:"page"`
	PageSize    int64       `json:"page_size"`
	SortFields  []sortField `json:"sort_fields"`
}

type registerStatsRequest struct {
	CreateBegin string      `json:"create_begin"`
	CreateEnd   string      `json:"create_end"`
}

type getUserListResponse struct {
	UserInfoList []userInfo `json:"user_info_list"`
	TotalCount   int64      `json:"total_count"`
}

type userInfo struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	Age       int8   `json:"age"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
}

type getUserDetailRequest struct {
	Email  string `json:"email"`
	UserId int64  `json:"user_id"`
}

type getUserDetailResponse struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	Age       int64  `json:"age"`
	Phone     string `json:"phone"`
	Gender    string `json:"gender"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
}

type userDetailRole struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type modifyUserDetailRequest struct {
	Email       string  `json:"email"`
	Name        string  `json:"name"`
	Status      int8    `json:"status"`
	AddRoles    []int64 `json:"add_roles"`
	DeleteRoles []int64 `json:"delete_roles"`
}

type sortField struct {
	Field string `json:"field"` // 必须和字段名一致
	Order string `json:"order"`
}

type userInfoCsvRow struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
	Age       string `json:"age"`
	Gender    string `json:"gender"`
	CreatedAt string `json:"created_at"`
}

type userRegisterInfoResponse struct {
	Series []registerInfo `json:"series"`
}

type registerInfo struct {
	Date string `json:"date"`
	Cnt  int64  `json:"cnt"`
}

type addUserRoleRequest struct {
	UserId int64 `json:"user_id"`
	RoleId int64 `json:"role_id"`
	Expire int64 `json:"expire"`
}

type deleteUserRoleRequest struct {
	UserId int64 `json:"user_id"`
	RoleId int64 `json:"role_id"`
}

type getUserOperationListRequest struct {
	UserId int64 `json:"user_id"`
}

type getUserOperationListResponse struct {
	List []userOperation `json:"list"`
}

type userOperation struct {
	Operation  string `json:"operation"`
	OperatedAt string `json:"operated_at"`
}
