// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package summary

import (
	"context"
	"sort"
	"strings"

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

	periods, err := l.svcCtx.Repos.Summary.ListTagPeriods(l.ctx, authUser.UserID, req.Start, req.End)
	if err != nil {
		l.Errorf("list tag periods failed: %v", err)
		return nil, errorx.WrapDBQuery("查询标签趋势失败", err)
	}

	// 按月分组统计标签
	monthMap := make(map[string]map[string]int64)
	for _, p := range periods {
		if monthMap[p.Month] == nil {
			monthMap[p.Month] = make(map[string]int64)
		}
		for _, tag := range strings.Split(p.Tags, ",") {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				monthMap[p.Month][tag]++
			}
		}
	}

	// 转换为有序结果
	var months []types.MonthTagStats
	for month, tagCounts := range monthMap {
		var tags []types.TagCount
		for t, c := range tagCounts {
			tags = append(tags, types.TagCount{Tag: t, Count: c})
		}
		sort.Slice(tags, func(i, j int) bool { return tags[i].Count > tags[j].Count })
		months = append(months, types.MonthTagStats{Month: month, Tags: tags})
	}
	sort.Slice(months, func(i, j int) bool { return months[i].Month < months[j].Month })

	return &types.SummaryTagTrendResp{Months: months}, nil
}
