package config

import (
	"time"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	MySQLConf struct {
		DSN string
	}
	RedisConf     RedisConf
	EmailConf     EmailConf
	IP2RegionConf IP2RegionConf
	JWTConf       JWTConf
	AIConf        AIConf
}

type AIConf struct {
	Endpoint string // DeepSeek API endpoint
	APIKey   string
	Model    string // 默认 deepseek-chat
}

type IP2RegionConf struct {
	V4 string
	V6 string
}

type EmailConf struct {
	From     string
	Password string
	SMTPHost string
	SMTPPort int
}

type JWTConf struct {
	Secret         string
	ExpireS        int64 // access token 过期时间，单位：秒，默认 900（15min）
	RefreshExpireS int64 // refresh token 过期时间，单位：秒，默认 604800（7d）
}

// RefreshExpireDuration 返回 refresh token 的过期时间 Duration
func (c JWTConf) RefreshExpireDuration() time.Duration {
	return time.Duration(c.RefreshExpireS) * time.Second
}

type RedisConf struct {
	Addr     string
	Password string
	DB       int
}
