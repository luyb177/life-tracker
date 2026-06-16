package constvar

import (
	"time"
)

// Context keys
type contextKey string

const (
	AuthUserKey   contextKey = "auth_user"
	IPLocationKey contextKey = "ip_location"
)

// 渠道
const (
	ChannelEmail int32 = 1
	ChannelPhone int32 = 2
)

// 验证码用途
const (
	PurposeRegistration  int32 = 1
	PurposePasswordReset int32 = 2
)

const VerifyCodeExpire = 5 * time.Minute

// 分页
const (
	DefaultPageSize = 20
)

// TimeLocation 数据库时区，与 DSN 中 loc 参数保持一致
var TimeLocation, _ = time.LoadLocation("Asia/Shanghai")
