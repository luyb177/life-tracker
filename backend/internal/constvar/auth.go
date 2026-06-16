package constvar

import "time"

// Context keys
type contextKey string

const (
	AuthUserKey contextKey = "auth_user"
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
