// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package summary

import (
	"context"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteSummaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSummaryLogic {
	return &DeleteSummaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSummaryLogic) DeleteSummary(req *types.DeleteSummaryReq) (*types.Response, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	s, err := l.svcCtx.Repos.Summary.FindByID(l.ctx, req.ID)
	if err != nil {
		l.Errorf("find summary failed: %v", err)
		return nil, errorx.WrapDBQuery("查询总结失败", err)
	}
	if s == nil {
		return nil, errorx.ErrNotFound
	}
	if s.UserID != authUser.UserID {
		return nil, errorx.ErrForbidden
	}

	if err := l.svcCtx.Repos.Transaction(func(tx *gorm.DB) error {
		if err := l.svcCtx.Repos.Tag.DeleteBySummaryID(l.ctx, req.ID, tx); err != nil {
			l.Errorf("delete summary tags failed: %v", err)
			return errorx.WrapDBDelete("删除标签关联失败", err)
		}

		if err := l.svcCtx.Repos.Summary.Delete(l.ctx, req.ID, tx); err != nil {
			l.Errorf("delete summary failed: %v", err)
			return errorx.WrapDBDelete("删除总结失败", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &types.Response{}, nil
}
