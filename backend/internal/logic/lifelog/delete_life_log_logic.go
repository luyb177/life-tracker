// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package lifelog

import (
	"context"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLifeLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除生活记录
func NewDeleteLifeLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLifeLogLogic {
	return &DeleteLifeLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLifeLogLogic) DeleteLifeLog(req *types.DeleteLifeLogReq) (resp *types.Response, err error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	existing, err := l.svcCtx.Repos.LifeLog.FindByID(l.ctx, req.ID)
	if err != nil {
		l.Errorf("find life log failed: %v", err)
		return nil, errorx.WrapDBQuery("查询生活记录失败", err)
	}
	if existing == nil {
		return nil, errorx.ErrNotFound
	}
	if existing.UserID != authUser.UserID {
		return nil, errorx.ErrForbidden
	}

	if err := l.svcCtx.Repos.LifeLog.Delete(l.ctx, req.ID); err != nil {
		l.Errorf("delete life log failed: %v", err)
		return nil, errorx.WrapDBDelete("删除生活记录失败", err)
	}

	return &types.Response{}, nil
}
