package authority

import (
	"fsbm/util"
	"time"
)

type getRoleListResponse struct {
	List []roleInfo `json:"list"`
}

type roleInfo struct {
	Role string `json:"role"`
	Id   int64  `json:"id"`
}

type getUserRoleListRequest struct {
	UserId int64  `json:"user_id"`
	Email  string `json:"email"`
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
	RoleId     int64  `json:"role_id"`
	Email      string `json:"email"`
	Reason     string `json:"reason"`
	Expiration int64  `json:"expiration"`
}

type applyRoleListRequest struct {
	User            string   `json:"user"`
	Role            []string `json:"role"`
	Reviewer        string   `json:"reviewer"`
	Status          []int8   `json:"status"`
	ApplyBeginTime  string   `json:"apply_begin_time"`
	ApplyEndTime    string   `json:"apply_end_time"`
	ReviewBeginTime string   `json:"review_begin_time"`
	ReviewEndTime   string   `json:"review_end_time"`
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
	CreatedAt    time.Time
}

type applyRoleCsvRow struct {
	User         string `json:"user"`
	Role         string `json:"role"`
	Reason       string `json:"reason"`
	Status       string `json:"status"`
	Reviewer     string `json:"reviewer"`
	ReviewReason string `json:"review_reason"`
	ReviewAt     string `json:"review_at"`
	CreatedAt    string `json:"created_at"`
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
	CreatedAt    string `json:"created_at"`
}

type reviewApplyRoleRequest struct {
	Id       int64  `json:"id"`
	Review   int8   `json:"review"`
	Reason   string `json:"reason"`
	Reviewer string `json:"reviewer"`
}
