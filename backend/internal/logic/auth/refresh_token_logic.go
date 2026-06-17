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

	// 2. 必须是 refresh 类型，且必须有 JTI（拒绝旧版无 JTI 的 token）
	if claims.TokenType != jwtx.TokenTypeRefresh {
		return nil, errorx.ErrTokenInvalid
	}
	if claims.ID == "" {
		l.Errorf("refresh token missing jti, token was generated before JTI rotation was implemented")
		return nil, errorx.ErrTokenInvalid
	}

	// 3. 签发新的 access + refresh token（每次刷新都生成新 JTI）
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

	// 4. 解析新 refresh token 获取 JTI
	newClaims, err := l.svcCtx.JWTHandler.ParseJWTToken(refreshToken)
	if err != nil {
		l.Errorf("parse new refresh token failed: %v", err)
		return nil, errorx.WrapInternal("令牌解析失败", err)
	}

	// 5. 原子性轮换：验证旧 JTI 并替换为新 JTI
	//    Lua 脚本保证 验证+替换 一步完成，旧 token 立即失效且不可重放
	rotated, err := l.svcCtx.Repos.Token.Rotate(l.ctx, claims.UserID, claims.ID, newClaims.ID, l.svcCtx.Config.JWTConf.RefreshExpireDuration())
	if err != nil {
		l.Errorf("rotate refresh token failed: %v", err)
		return nil, errorx.WrapRedisSet("令牌刷新失败", err)
	}
	if !rotated {
		l.Errorf("refresh token jti mismatch: user=%d", claims.UserID)
		return nil, errorx.ErrTokenInvalid
	}

	return &types.RefreshTokenResp{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}
