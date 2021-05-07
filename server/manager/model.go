package manager

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
