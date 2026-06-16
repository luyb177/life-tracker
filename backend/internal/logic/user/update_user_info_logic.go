// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfoReq) (*types.Response, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["name"] = req.Username
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if len(updates) == 0 {
		return nil, errorx.WrapBadRequest("没有可更新的字段", nil)
	}

	if err := l.svcCtx.Repos.User.Update(l.ctx, authUser.UserID, updates); err != nil {
		l.Errorf("update user failed: %v", err)
		return nil, errorx.WrapDBUpdate("更新用户信息失败", err)
	}

	return &types.Response{}, nil
}
