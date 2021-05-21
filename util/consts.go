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
	ApplyROleOrderHasReviewed = 3002 // 工单已被审批
	ShopNotFound              = 4001 // 店铺信息不存在
	ShopNotUniq               = 4002 // 店铺查询结果不唯一
	ShopNotBelongAdjuster     = 4003 // 店铺不属于操作者
	AbnormalError             = 5001 // 异常错误（预料外错误统一用这个）
	// redis key template
	UserLoginTemplate            = "%s:login"
	UserVerificationCodeTemplate = "%s:verification_code"
	DashboardRecordCnt           = "%s:dashboard_record_cnt"
	DashboardAlarmCnt            = "%s:dashboard_alarm_cnt"
	DashboardLatestRecord        = "%s:dashboard_latest_record"
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
	// 用户操作记录
	UserOperation_AddShop       = "注册店铺[%d]%s"
	UserOperation_AddMonitor    = "为店铺[%d]注册监控[%d]%s"
	UserOperation_ModifyProfile = "修改个人信息,%s -> %s"
	UserOperation_DeleteUser    = "删除账户"
	UserOPeration_Register      = "注册账户"
	// 短信发送平台配置
	PhoneMessageUrl            = "http://gw.api.taobao.com/router/rest"
	PhoneMessageMethod         = "alibaba.aliqin.fc.sms.num.send"
	PhoneMessageTemplate       = "" // 模版id
	PhoneMessageSecret         = "" // app secret
	PhoneMessageAppKey         = "" // app key
	PhoneMessageSignName       = "" // 签名
	PhoneMessageV2TemplateCode = "SMS_217415565"
	PhoneMessageV2Url          = "http://www.weikebaijia.net/ylxxt/sms_message_servlet_action"
	PhoneMessageV2Action       = "send_template_sms_message"
	// 微信公众号消息发送配置
	WxMessageTemplate = "f8nZklEIKNvW1ziAnzIPs-Lc3s0hGfP-aIqRa9n4Llc"
	WxMessageUrl      = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=45_RGW2cXRcVh5-pBNmSp6DTPN4qWlRjngAub3PAE4Z68pA_2CkIAuQdPIht5rFyML7ydn1iBdveoquqdFFnABDf_WeT7Sc26hSaUfxMnlx4aCVjgmNdz7SmsMK2OhTR-91J7MRpTDe4k23tSmsULTfAFAZQY"
)

// 运行时初始化常数
var (
	Role_NormalUserId  int64 // 普通用户角色id
	Role_AdminId       int64 // 管理员角色id
	Role_supervisionId int64 // 监管角色id
	Role_ShopOwnerId   int64 // 店铺管理员角色id
)
