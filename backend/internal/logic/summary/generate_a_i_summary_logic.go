// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package summary

import (
	"context"
	"errors"
	"fmt"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/repo/summary"
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

	// TODO: 接入 LLM 真实生成
	label := periodTypeLabel(req.PeriodType)
	content := fmt.Sprintf("这是您的%s（周期：%s）。AI 总结功能即将上线。", label, req.PeriodStart)
	suggestion := "功能开发中，敬请期待。"

	s := &summary.Summary{
		UserID:            authUser.UserID,
		PeriodType:        req.PeriodType,
		PeriodStart:       req.PeriodStart,
		PeriodEnd:         req.PeriodStart,
		Source:            constvar.SummarySourceAI,
		SummaryContent:    content,
		SuggestionContent: suggestion,
	}

	existing, err := l.svcCtx.Repos.Summary.FindByPeriod(l.ctx, authUser.UserID, req.PeriodType, req.PeriodStart)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		l.Errorf("find existing summary failed: %v", err)
		return nil, errorx.WrapDBQuery("查询已有总结失败", err)
	}

	if existing != nil && existing.Source == constvar.SummarySourceAI {
		updates := map[string]interface{}{
			"summary_content":    content,
			"suggestion_content": suggestion,
		}
		if err := l.svcCtx.Repos.Summary.Update(l.ctx, existing.ID, updates); err != nil {
			l.Errorf("update ai summary failed: %v", err)
			return nil, errorx.WrapDBUpdate("更新AI总结失败", err)
		}
		s.ID = existing.ID
		s.CreatedAt = existing.CreatedAt
	} else {
		if err := l.svcCtx.Repos.Summary.Create(l.ctx, s); err != nil {
			l.Errorf("create ai summary failed: %v", err)
			return nil, errorx.WrapDBInsert("创建AI总结失败", err)
		}
	}

	return &types.SummaryInfo{
		ID:                s.ID,
		PeriodType:        s.PeriodType,
		PeriodStart:       s.PeriodStart,
		PeriodEnd:         s.PeriodEnd,
		Source:            s.Source,
		SummaryContent:    s.SummaryContent,
		SuggestionContent: s.SuggestionContent,
		CreatedAt:         s.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         s.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
