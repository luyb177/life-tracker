// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"
	"strings"
	"time"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
	"github.com/luyb177/life-tracker/backend/internal/pkg/code"
	"github.com/luyb177/life-tracker/backend/internal/pkg/email"
	"github.com/luyb177/life-tracker/backend/internal/repo/verify"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const verifyCodeCooldown = 60 * time.Second

type SendVerificationCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendVerificationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendVerificationCodeLogic {
	return &SendVerificationCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendVerificationCodeLogic) SendVerificationCode(req *types.SendVerificationCodeReq) (*types.Response, error) {
	if !validVerifyPurpose(req.Purpose) {
		return nil, errorx.WrapBadRequest("无效的验证码用途", nil)
	}

	switch req.Channel {
	case constvar.ChannelEmail:
		return l.sendByEmail(req.Target, req.Purpose)
	default:
		return nil, errorx.WrapBadRequest("无效的验证码渠道", nil)
	}
}

func (l *SendVerificationCodeLogic) sendByEmail(target string, purpose int32) (*types.Response, error) {
	if strings.TrimSpace(target) == "" {
		return nil, errorx.WrapBadRequest("邮箱地址不能为空", nil)
	}
	if len(target) > 254 {
		return nil, errorx.WrapBadRequest("邮箱地址过长", nil)
	}
	target = email.CanonicalEmail(target)
	if !email.IsValidEmail(target) {
		return nil, errorx.WrapBadRequest("无效的邮箱地址", nil)
	}

	// 限频：60 秒内同目标同用途不允许重复发送
	meta := &verify.Meta{Target: target, Channel: constvar.ChannelEmail, Purpose: purpose}
	if ttl, err := l.svcCtx.Repos.Verify.CodeTTL(l.ctx, meta); err == nil && ttl > 0 {
		remaining := constvar.VerifyCodeExpire - ttl
		if remaining < verifyCodeCooldown {
			return nil, errorx.WrapBadRequest("验证码发送过于频繁，请稍后再试", nil)
		}
	}

	emailCode := code.EmailCode()

	if err := l.svcCtx.Repos.Verify.SetCode(l.ctx, meta, emailCode, constvar.VerifyCodeExpire); err != nil {
		l.Errorf("set verify code failed: %v", err)
		return nil, errorx.WrapInternal("验证码存储失败", err)
	}

	go func() {
		if err := l.svcCtx.EmailSender.SendVerifyCode(context.Background(), target, emailCode, int(constvar.VerifyCodeExpire.Minutes())); err != nil {
			l.Errorf("send verify code email failed: %v", err)
		}
	}()

	return &types.Response{}, nil
}

func validVerifyPurpose(purpose int32) bool {
	switch purpose {
	case constvar.PurposeRegistration, constvar.PurposePasswordReset:
		return true
	default:
		return false
	}
}
