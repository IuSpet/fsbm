package userAccount

type userCommonRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	Gender     int8   `json:"gender"`
	Age        int8   `json:"age"`
	VerifyCode string `json:"verify_code"`
}

type loginResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type applyRoleRequest struct {
	Email      string  `json:"email"`
	RoleIDList []int64 `json:"role_id_list"`
}

type getUserProfileRequest struct {
	Email string `json:"email"`
}

type getUserProfileResponse struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Gender    int8   `json:"gender"`
	Age       int8   `json:"age"`
	CreatedAt string `json:"created_at"`
}

type deleteUserRequest struct {
	Email      string `json:"email"`
	VerifyCode string `json:"verify_code"`
}

type getUserRolesResponse struct {
	Roles []string `json:"roles"`
}
