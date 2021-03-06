package userAccount

type userCommonRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	VerifyCode string `json:"verify_code"`
}

type applyRoleRequest struct {
	Email      string  `json:"email"`
	RoleIDList []int64 `json:"role_id_list"`
}
