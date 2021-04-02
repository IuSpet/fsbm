package util

const (
	// 所有返回status含义
	ParamError              = 1  // 请求参数错误
	DbError                 = 2  // 数据库错误，包括mysql和redis
	RepeatedEmailAddr       = 3  // 已被使用过的邮箱
	IllegalEmailAddr        = 4  // 邮箱格式错误
	UserNotExist            = 5  // 用户不存在
	InvalidPassword         = 6  // 密码错误
	InvalidVerificationCode = 7  // 验证码错误
	UserNotLogin            = 8  // 用户未登录
	AuthenticationFail      = 9  // 用户鉴权失败
	UserDeleted             = 10 // 用户已被删除
	EmailSendError          = 11 // 邮件发送失败
	AvatarNotExist          = 12 // 没有设置头像
	// redis key template
	UserLoginTemplate            = "%s:login_in"
	UserVerificationCodeTemplate = "%s:verification_code"
	// 密码加盐
	Salt = "fsbmpwd"
	// time format
	YMD    = "2006-01-02"
	YMDHMS = "2006-01-02 15:04:05"
	H5FMT  = "2006-01-02T15:04"
)
