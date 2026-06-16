package svc

import (
	"time"

	"github.com/luyb177/life-tracker/backend/common/cache"
	"github.com/luyb177/life-tracker/backend/common/database"
	"github.com/luyb177/life-tracker/backend/common/jwtx"
	"github.com/luyb177/life-tracker/backend/common/mail"
	"github.com/luyb177/life-tracker/backend/internal/config"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/pkg/email"
	"github.com/luyb177/life-tracker/backend/internal/repo"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config        config.Config
	Mailer        *mail.Mailer
	RedisClient   *cache.RedisClient
	MySQLClient   *database.MySQLClient
	EmailSender   email.EmailSender
	Repos         *repo.Repositories
	JWTHandler    jwtx.Handler
	IPMiddleware  rest.Middleware
	JWTMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	m := mail.NewMailer(mail.EmailConfig{
		From:     c.EmailConf.From,
		Password: c.EmailConf.Password,
		SMTPHost: c.EmailConf.SMTPHost,
		SMTPPort: c.EmailConf.SMTPPort,
	})

	emailSender := email.NewEmailSender(m)

	rc, err := cache.NewRedisClient(c.RedisConf.Addr, c.RedisConf.Password, c.RedisConf.DB)
	if err != nil {
		panic(err)
	}

	mc, err := database.NewMySQLClient(c.MySQLConf.DSN)
	if err != nil {
		panic(err)
	}

	jwtHandler, err := jwtx.NewHandler(
		c.JWTConf.Secret,
		time.Duration(c.JWTConf.ExpireS)*time.Second,
		time.Duration(c.JWTConf.RefreshExpireS)*time.Second,
	)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:        c,
		Mailer:        m,
		RedisClient:   rc,
		MySQLClient:   mc,
		EmailSender:   emailSender,
		Repos:         repo.NewRepositories(rc.Client, mc.DB),
		JWTHandler:    jwtHandler,
		IPMiddleware:  middleware.NewIPMiddleware(c.IP2RegionConf).Handle,
		JWTMiddleware: middleware.NewJWTMiddleware(jwtHandler).Handle,
	}
}
