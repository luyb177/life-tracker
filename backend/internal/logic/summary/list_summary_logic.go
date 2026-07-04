// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package summary

import (
	"context"
	"time"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/pkg/pagetoken"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSummaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSummaryLogic {
	return &ListSummaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListSummaryLogic) ListSummary(req *types.ListSummaryReq) (*types.ListSummaryResp, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	limit := int(req.PageSize)
	if limit <= 0 || limit > 50 {
		limit = constvar.DefaultPageSize
	}

	// 解码游标
	var cursorID uint64
	var cursorTime time.Time
	if req.PageToken != "" {
		var pt types.PageToken
		if err := pagetoken.Decode(req.PageToken, constvar.SummaryPageTokenPrefix, &pt); err != nil {
			return nil, errorx.WrapBadRequest("分页参数无效", err)
		}
		cursorID = pt.ID
		if pt.CreatedAt != "" {
			var parseErr error
			cursorTime, parseErr = time.ParseInLocation(time.DateTime, pt.CreatedAt, constvar.TimeLocation)
			if parseErr != nil {
				return nil, errorx.WrapBadRequest("分页参数无效", parseErr)
			}
		}
	}

	// 多查一条判断 HasMore
	list, err := l.svcCtx.Repos.Summary.ListByUser(l.ctx, authUser.UserID, req.PeriodType, cursorID, cursorTime, limit+1)
	if err != nil {
		l.Errorf("list summary failed: %v", err)
		return nil, errorx.WrapDBQuery("查询总结列表失败", err)
	}

	hasMore := len(list) > limit
	if hasMore {
		list = list[:limit]
	}

	// 批量填充标签
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

	var nextToken string
	if hasMore && len(items) > 0 {
		last := list[len(items)-1]
		nextPT := types.PageToken{
			ID:        last.ID,
			CreatedAt: last.CreatedAt.In(constvar.TimeLocation).Format(time.DateTime),
		}
		nextToken, _ = pagetoken.Encode(constvar.SummaryPageTokenPrefix, &nextPT)
	}

	return &types.ListSummaryResp{
		List:      items,
		PageToken: nextToken,
		HasMore:   hasMore,
	}, nil
}
