package authority

import "fsbm/util"

type getRoleListResponse struct {
	List []roleInfo `json:"list"`
}

type roleInfo struct {
	Role string `json:"role"`
	Id   int64  `json:"id"`
}

type getUserRoleListRequest struct {
	Email string `json:"email"`
}

type getUserRoleListResponse struct {
	ActiveRoles  []roleInfo `json:"active_roles"`
	ExpiredRoles []roleInfo `json:"expired_roles"`
}

type userRoleStatusInfo struct {
	Role   string
	RoleId int64
	Status int8
}

type applyRoleRequest struct {
	RoleId int64  `json:"role_id"`
	Email  string `json:"email"`
	Reason string `json:"reason"`
}

type applyRoleListRequest struct {
	User      string `json:"user"`
	Role      string `json:"role"`
	Reviewer  string `json:"reviewer"`
	Status    []int8 `json:"status"`
	BeginDate string `json:"begin_date"`
	EndDate   string `json:"end_date"`
	util.ListReqField
}

type applyRoleListResponse struct {
	List     []applyRoleOrder `json:"list"`
	TotalCnt int64            `json:"total_cnt"`
}

type applyRoleRow struct {
	Id           int64
	User         string
	Role         string
	Reason       string
	Status       int8
	Reviewer     string
	ReviewReason string
	ReviewAt     int64
}

type applyRoleOrder struct {
	Id           int64  `json:"id"`
	User         string `json:"user"`
	Role         string `json:"role"`
	Reason       string `json:"reason"`
	Status       string `json:"status"`
	Reviewer     string `json:"reviewer"`
	ReviewReason string `json:"review_reason"`
	ReviewAt     string `json:"review_at"`
}
