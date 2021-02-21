package util

const (
	// 所有返回status含义
	ParamError              = 1
	DbError                 = 2
	RepeatedEmailAddr       = 3
	IllegalEmailAddr        = 4
	UserNotExist            = 5
	InvalidPassword         = 6
	InvalidVerificationCode = 7
	UserNotLogin            = 8
	AuthenticationFail      = 9
	UserDeleted             = 10
	EmailSendError          = 11
	// redis key template
	UserLoginTemplate                 = "%s:login_in"
	UserLoginVerificationCodeTemplate = "%s:login_verification_code"
	// 密码加盐
	Salt = "fsbmpwd"
)
