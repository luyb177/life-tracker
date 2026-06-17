// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package summary

import (
	"context"
	"strings"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/repo/summary"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSummaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSummaryLogic {
	return &CreateSummaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSummaryLogic) CreateSummary(req *types.CreateSummaryReq) (*types.IDResponse, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	if strings.TrimSpace(req.SummaryContent) == "" {
		return nil, errorx.WrapBadRequest("总结内容不能为空", nil)
	}

	if !validPeriodType(req.PeriodType) {
		return nil, errorx.WrapBadRequest("无效的周期类型", nil)
	}

	s := &summary.Summary{
		UserID:            authUser.UserID,
		PeriodType:        req.PeriodType,
		PeriodStart:       req.PeriodStart,
		PeriodEnd:         req.PeriodEnd,
		Source:            constvar.SummarySourceUser,
		SummaryContent:    req.SummaryContent,
		SuggestionContent: req.SuggestionContent,
		Location:          middleware.FullLocation(middleware.GetIPLocation(l.ctx)),
	}

	if err := l.svcCtx.Repos.Summary.Create(l.ctx, s); err != nil {
		l.Errorf("create summary failed: %v", err)
		return nil, errorx.WrapDBInsert("创建总结失败", err)
	}

	return &types.IDResponse{
		ID: s.ID,
	}, nil
}
