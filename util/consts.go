package util

import "time"

type NoticeLevel string

const (
	// 所有返回status含义
	ParamError                = 1001 // 请求参数错误
	DbError                   = 1002 // 数据库错误，包括mysql和redis
	InvalidVerificationCode   = 1003 // 验证码错误
	EmailSendError            = 1004 // 邮件发送失败
	SaveImgError              = 1005 // 保存图片到本地失败
	RepeatedEmailAddr         = 2001 // 已被使用过的邮箱
	IllegalEmailAddr          = 2002 // 邮箱格式错误
	UserNotExist              = 2003 // 用户不存在
	InvalidPassword           = 2004 // 密码错误
	UserNotLogin              = 2005 // 用户未登录
	AuthenticationFail        = 2006 // 用户鉴权失败
	UserDeleted               = 2007 // 用户已被删除
	AvatarNotExist            = 2008 // 没有设置头像
	ApplyRoleOrderNotExist    = 3001 // 角色申请工单不存在
	ApplyROleOrderHasReviewed = 3002
	AbnormalError             = 4001 // 异常错误（预料外错误统一用这个）
	// redis key template
	UserLoginTemplate            = "%s:login"
	UserVerificationCodeTemplate = "%s:verification_code"
	// 密码加盐
	Salt = "fsbmpwd"
	// time format
	YMD    = "2006-01-02"
	YMDHMS = "2006-01-02 15:04:05"
	H5FMT  = "2006-01-02T15:04"
	// 报警时间配置
	InfoNoticeInterval  = time.Hour
	WarnNoticeInterval  = 15 * time.Minute
	ErrorNoticeInterval = time.Minute
	// 报警等级枚举
	InfoNotice  NoticeLevel = "_info_notice"
	WarnNotice  NoticeLevel = "_warn_notice"
	ErrorNotice NoticeLevel = "_error_notice"
)
