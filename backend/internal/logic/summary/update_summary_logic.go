// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package summary

import (
	"context"
	"strings"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
	if req.Tags != "" {
		updates["tags"] = strings.TrimSpace(req.Tags)
	}
	if len(updates) == 0 {
		return nil, errorx.WrapBadRequest("没有可更新的字段", nil)
	}

	if err := l.svcCtx.Repos.Summary.Update(l.ctx, req.ID, updates); err != nil {
		l.Errorf("update summary failed: %v", err)
		return nil, errorx.WrapDBUpdate("更新总结失败", err)
	}

	return &types.Response{}, nil
}
