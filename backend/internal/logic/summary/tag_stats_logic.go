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

type TagStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTagStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagStatsLogic {
	return &TagStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagStatsLogic) TagStats(req *types.SummaryTagStatsReq) (*types.SummaryTagStatsResp, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	tags, err := l.svcCtx.Repos.Summary.ListTagsByDateRange(l.ctx, authUser.UserID, req.Start, req.End)
	if err != nil {
		l.Errorf("list tags failed: %v", err)
		return nil, errorx.WrapDBQuery("查询标签统计失败", err)
	}

	// 按逗号拆分标签并统计频次
	countMap := make(map[string]int64)
	for _, tagStr := range tags {
		for _, tag := range strings.Split(tagStr, ",") {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				countMap[tag]++
			}
		}
	}

	// 按频次降序排列
	type kv struct {
		tag   string
		count int64
	}
	var sorted []kv
	for t, c := range countMap {
		sorted = append(sorted, kv{t, c})
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].count > sorted[j].count })

	result := make([]types.TagCount, 0, len(sorted))
	for _, kv := range sorted {
		result = append(result, types.TagCount{Tag: kv.tag, Count: kv.count})
	}

	return &types.SummaryTagStatsResp{Tags: result}, nil
}
