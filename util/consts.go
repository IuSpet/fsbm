package util

const (
	// 所有返回status含义
	ParamError        = 1
	DbError           = 2
	RepeatedEmailAddr = 3
	IllegalEmailAddr  = 4
	UserNotExist      = 5
	InvalidPassword   = 6
	// redis key template
	UserLoginTemplate = "%s:login_in"
	// 密码加盐
	Salt = "fsbmpwd"
)
