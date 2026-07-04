// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package summary

import (
	"context"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RangeSummaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRangeSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RangeSummaryLogic {
	return &RangeSummaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RangeSummaryLogic) RangeSummary(req *types.SummaryRangeReq) (*types.SummaryRangeResp, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	if req.Start == "" || req.End == "" {
		return nil, errorx.WrapBadRequest("start 和 end 不能为空", nil)
	}

	list, err := l.svcCtx.Repos.Summary.FindByPeriodRange(l.ctx, authUser.UserID, req.PeriodType, req.Start, req.End)
	if err != nil {
		l.Errorf("query range summary failed: %v", err)
		return nil, errorx.WrapDBQuery("查询总结失败", err)
	}

	ids := make([]uint64, 0, len(list))
	for _, s := range list {
		ids = append(ids, s.ID)
	}
	tagMap, err := batchFillSummaryTags(l.ctx, l.svcCtx, ids)
	if err != nil {
		l.Errorf("batch fill tags failed: %v", err)
		return nil, errorx.WrapDBQuery("查询标签失败", err)
	}

	items := make([]types.SummaryInfo, 0, len(list))
	for _, s := range list {
		tagInfos := tagMap[s.ID]
		if tagInfos == nil {
			tagInfos = []types.TagInfo{}
		}
		items = append(items, summaryToInfo(s, tagInfos))
	}

	return &types.SummaryRangeResp{List: items}, nil
}
