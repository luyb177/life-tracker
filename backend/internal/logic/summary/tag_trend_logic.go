// Code scaffolded by goctl. Safe to edit.
package summary

import (
	"context"
	"sort"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagTrendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTagTrendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagTrendLogic {
	return &TagTrendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagTrendLogic) TagTrend(req *types.SummaryTagStatsReq) (*types.SummaryTagTrendResp, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	freqs, err := l.svcCtx.Repos.Tag.ListTagMonthFrequencies(l.ctx, authUser.UserID, req.Start, req.End)
	if err != nil {
		l.Errorf("list tag month frequencies failed: %v", err)
		return nil, errorx.WrapDBQuery("查询标签趋势失败", err)
	}

	// Group by month
	monthMap := make(map[string][]types.TagCount)
	for _, f := range freqs {
		monthMap[f.Month] = append(monthMap[f.Month], types.TagCount{Tag: f.Tag, Count: f.Count})
	}

	var months []types.MonthTagStats
	for month, tags := range monthMap {
		sort.Slice(tags, func(i, j int) bool { return tags[i].Count > tags[j].Count })
		months = append(months, types.MonthTagStats{Month: month, Tags: tags})
	}
	sort.Slice(months, func(i, j int) bool { return months[i].Month < months[j].Month })

	return &types.SummaryTagTrendResp{Months: months}, nil
}
