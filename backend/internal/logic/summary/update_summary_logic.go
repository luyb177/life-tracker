// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package summary

import (
	"context"
	"strings"
	"time"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateSummaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSummaryLogic {
	return &UpdateSummaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSummaryLogic) UpdateSummary(req *types.UpdateSummaryReq) (*types.Response, error) {
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

	updates := make(map[string]interface{})
	if req.SummaryContent != "" {
		updates["summary_content"] = req.SummaryContent
	}
	if req.SuggestionContent != "" {
		updates["suggestion_content"] = req.SuggestionContent
	}
	if req.Title != "" {
		updates["title"] = strings.TrimSpace(req.Title)
	}
	if len(updates) == 0 && req.Tags == nil {
		return nil, errorx.WrapBadRequest("没有可更新的字段", nil)
	}
	updates["last_updated_by"] = authUser.UserID
	updates["last_updated_at"] = time.Now().In(constvar.TimeLocation)

	if err := l.svcCtx.Repos.Transaction(func(tx *gorm.DB) error {
		// 标签替换
		if req.Tags != nil {
			if err := l.svcCtx.Repos.Tag.DeleteBySummaryID(l.ctx, req.ID, tx); err != nil {
				l.Errorf("delete old tags failed: %v", err)
				return errorx.WrapDBDelete("删除旧标签关联失败", err)
			}
			if err := resolveSummaryTags(l.ctx, l.svcCtx, req.ID, req.Tags, tx); err != nil {
				return err
			}
		}

		if err := l.svcCtx.Repos.Summary.Update(l.ctx, req.ID, updates, tx); err != nil {
			l.Errorf("update summary failed: %v", err)
			return errorx.WrapDBUpdate("更新总结失败", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &types.Response{}, nil
}
