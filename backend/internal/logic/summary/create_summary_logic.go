// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package summary

import (
	"context"
	"fmt"
	"strings"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/repo/summary"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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

	startAt, endAt, err := normalizePeriodRange(req.PeriodType, req.PeriodStart, req.PeriodEnd)
	if err != nil {
		return nil, errorx.WrapBadRequest(fmt.Sprintf("周期参数无效：%v；period_start / period_end 期望格式: %s", err, periodStartHint(req.PeriodType)), err)
	}
	normalizedStart := startAt.Format("2006-01-02")
	normalizedEnd := endAt.Format("2006-01-02")

	// 同周期仅允许一条用户记录
	exists, err := l.svcCtx.Repos.Summary.ExistsByPeriodAndSource(l.ctx, authUser.UserID, req.PeriodType, normalizedStart, constvar.SummarySourceUser)
	if err != nil {
		l.Errorf("check summary exists failed: %v", err)
		return nil, errorx.WrapDBQuery("查询已有记录失败", err)
	}
	if exists {
		return nil, errorx.WrapBadRequest(fmt.Sprintf("该%s已存在，不允许重复创建", periodTypeLabel(req.PeriodType)), nil)
	}

	s := &summary.Summary{
		UserID:            authUser.UserID,
		PeriodType:        req.PeriodType,
		PeriodStart:       normalizedStart,
		PeriodEnd:         normalizedEnd,
		Source:            constvar.SummarySourceUser,
		SummaryContent:    req.SummaryContent,
		SuggestionContent: req.SuggestionContent,
		Title:             strings.TrimSpace(req.Title),
		Location:          middleware.FullLocation(middleware.GetIPLocation(l.ctx)),
	}

	if err := l.svcCtx.Repos.Transaction(func(tx *gorm.DB) error {
		if err := l.svcCtx.Repos.Summary.Create(l.ctx, s, tx); err != nil {
			l.Errorf("create summary failed: %v", err)
			return errorx.WrapDBInsert("创建总结失败", err)
		}

		if err := resolveSummaryTags(l.ctx, l.svcCtx, s.ID, req.Tags, tx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &types.IDResponse{ID: s.ID}, nil
}
