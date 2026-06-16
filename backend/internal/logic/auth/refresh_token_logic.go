// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/common/jwtx"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenReq) (*types.RefreshTokenResp, error) {
	// 1. 解析 refresh token
	claims, err := l.svcCtx.JWTHandler.ParseJWTToken(req.RefreshToken)
	if err != nil {
		l.Errorf("parse refresh token failed: %v", err)
		return nil, errorx.ErrTokenInvalid
	}

	// 2. 必须是 refresh 类型
	if claims.TokenType != jwtx.TokenTypeRefresh {
		return nil, errorx.ErrTokenInvalid
	}

	// 3. 验证 JTI 是否仍在 Redis 中（未被刷新覆盖 / 未过期）
	valid, err := l.svcCtx.Repos.Token.Validate(l.ctx, claims.UserID, claims.JTI)
	if err != nil {
		l.Errorf("validate refresh token jti failed: %v", err)
		return nil, errorx.WrapRedisGet("令牌校验失败", err)
	}
	if !valid {
		return nil, errorx.ErrTokenInvalid
	}

	// 4. 签发新的 access + refresh token（每次刷新都生成新 JTI）
	accessToken, err := l.svcCtx.JWTHandler.SetAccessToken(jwtx.ClaimsParams{
		UserID: claims.UserID,
	})
	if err != nil {
		l.Errorf("generate access token failed: %v", err)
		return nil, errorx.WrapInternal("生成令牌失败", err)
	}

	refreshToken, err := l.svcCtx.JWTHandler.SetRefreshToken(jwtx.ClaimsParams{
		UserID: claims.UserID,
	})
	if err != nil {
		l.Errorf("generate refresh token failed: %v", err)
		return nil, errorx.WrapInternal("生成令牌失败", err)
	}

	// 5. 存储新 JTI，覆盖旧值 → 旧 refresh token 立即失效（防重放）
	newClaims, err := l.svcCtx.JWTHandler.ParseJWTToken(refreshToken)
	if err != nil {
		l.Errorf("parse new refresh token failed: %v", err)
		return nil, errorx.WrapInternal("令牌解析失败", err)
	}
	if err := l.svcCtx.Repos.Token.Store(l.ctx, claims.UserID, newClaims.JTI, l.svcCtx.Config.JWTConf.RefreshExpireDuration()); err != nil {
		l.Errorf("store new refresh token jti failed: %v", err)
		return nil, errorx.WrapRedisSet("存储令牌失败", err)
	}

	return &types.RefreshTokenResp{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}
