// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
	"github.com/luyb177/life-tracker/backend/internal/pkg/email"
	"github.com/luyb177/life-tracker/backend/internal/pkg/password"
	"github.com/luyb177/life-tracker/backend/internal/repo/user"
	"github.com/luyb177/life-tracker/backend/internal/repo/verify"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (*types.Response, error) {
	if resp, err := l.validReq(req); err != nil {
		return resp, err
	}

	target := email.CanonicalEmail(req.Target)

	switch req.Channel {
	case constvar.ChannelEmail:
		return l.registerByEmail(target, req)
	default:
		return nil, errorx.WrapBadRequest("暂仅支持邮箱注册", nil)
	}
}

func (l *RegisterLogic) registerByEmail(target string, req *types.RegisterReq) (*types.Response, error) {
	meta := &verify.Meta{
		Target:  target,
		Channel: constvar.ChannelEmail,
		Purpose: constvar.PurposeRegistration,
	}
	ok, err := l.svcCtx.Repos.Verify.VerifyCode(l.ctx, meta, req.Code)
	if err != nil {
		l.Errorf("verify code failed: %v", err)
		return nil, errorx.WrapInternal("验证码校验失败", err)
	}
	if !ok {
		return nil, errorx.WrapBadRequest("验证码错误或已过期", nil)
	}

	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		l.Errorf("hash password failed: %v", err)
		return nil, errorx.WrapInternal("密码加密失败", err)
	}

	u := &user.User{
		Name:     randomName(),
		Email:    target,
		Password: hashedPassword,
	}
	if err := l.svcCtx.Repos.User.Create(l.ctx, u); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errorx.WrapBadRequest("该邮箱已被注册", nil)
		}
		l.Errorf("create user failed: %v", err)
		return nil, errorx.WrapDBInsert("创建用户失败", err)
	}

	return &types.Response{}, nil
}

func (l *RegisterLogic) validReq(req *types.RegisterReq) (*types.Response, error) {
	target := strings.TrimSpace(req.Target)
	code := strings.TrimSpace(req.Code)
	pwd := req.Password

	switch {
	case target == "":
		return nil, errorx.WrapBadRequest("邮箱不能为空", nil)
	case len(target) > 254:
		return nil, errorx.WrapBadRequest("邮箱地址过长", nil)
	case !email.IsValidEmail(email.CanonicalEmail(target)):
		return nil, errorx.WrapBadRequest("无效的邮箱格式", nil)
	case code == "":
		return nil, errorx.WrapBadRequest("验证码不能为空", nil)
	case len(code) != 6:
		return nil, errorx.WrapBadRequest("验证码长度不正确", nil)
	case pwd == "":
		return nil, errorx.WrapBadRequest("密码不能为空", nil)
	case len(pwd) < 6:
		return nil, errorx.WrapBadRequest("密码长度不能少于6位", nil)
	case len(pwd) > 128:
		return nil, errorx.WrapBadRequest("密码长度不能超过128位", nil)
	}

	return nil, nil
}

func randomName() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return "用户_" + hex.EncodeToString(b)
}
