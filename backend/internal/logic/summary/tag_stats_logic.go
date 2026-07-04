// Code scaffolded by goctl. Safe to edit.
package summary

import (
	"context"

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

	freqs, err := l.svcCtx.Repos.Tag.ListTagFrequencies(l.ctx, authUser.UserID, req.Start, req.End)
	if err != nil {
		l.Errorf("list tag frequencies failed: %v", err)
		return nil, errorx.WrapDBQuery("查询标签统计失败", err)
	}

	result := make([]types.TagCount, 0, len(freqs))
	for _, f := range freqs {
		result = append(result, types.TagCount{Tag: f.Tag, Count: f.Count})
	}

	return &types.SummaryTagStatsResp{Tags: result}, nil
}
