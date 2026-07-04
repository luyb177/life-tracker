// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package summary

import (
	"context"
	"errors"
	"fmt"

	cronSummary "github.com/luyb177/life-tracker/backend/internal/logic/cron"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GenerateAISummaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateAISummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateAISummaryLogic {
	return &GenerateAISummaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateAISummaryLogic) GenerateAISummary(req *types.GenerateAISummaryReq) (*types.SummaryInfo, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	if !validPeriodType(req.PeriodType) {
		return nil, errorx.WrapBadRequest("无效的周期类型", nil)
	}

	periodStart, err := normalizePeriodStart(req.PeriodType, req.PeriodStart)
	if err != nil {
		return nil, errorx.WrapBadRequest(fmt.Sprintf("周期起始日期无效：%v；period_start 期望格式: %s", err, periodStartHint(req.PeriodType)), err)
	}
	periodStartKey := periodStart.Format("2006-01-02")

	// 调用真实 AI 总结流程
	if err := cronSummary.Run(l.ctx, l.svcCtx, req.PeriodType, authUser.UserID, periodStart); err != nil {
		l.Errorf("run ai summary failed: %v", err)
		return nil, errorx.WrapInternal("AI 总结生成失败", err)
	}

	// 查询刚生成的 AI 总结（按 source=AI 精确回查）
	s, err := l.svcCtx.Repos.Summary.FindByPeriodAndSource(l.ctx, authUser.UserID, req.PeriodType, periodStartKey, constvar.SummarySourceAI)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.ErrNotFound
		}
		l.Errorf("find created summary failed: %v", err)
		return nil, errorx.WrapDBQuery("查询总结失败", err)
	}
	if s == nil {
		return nil, errorx.ErrNotFound
	}

	// 查询标签
	tagMap, _ := batchFillSummaryTags(l.ctx, l.svcCtx, []uint64{s.ID})
	tagInfos := tagMap[s.ID]
	if tagInfos == nil {
		tagInfos = []types.TagInfo{}
	}

	return &types.SummaryInfo{
		ID:                s.ID,
		PeriodType:        s.PeriodType,
		PeriodStart:       s.PeriodStart,
		PeriodEnd:         s.PeriodEnd,
		Source:            s.Source,
		SummaryContent:    s.SummaryContent,
		SuggestionContent: s.SuggestionContent,
		Title:             s.Title,
		Tags:              tagInfos,
		Location:          s.Location,
		CreatedAt:         s.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         s.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
