// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package lifelog

import (
	"context"
	"time"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/pkg/pagetoken"
	lifelog "github.com/luyb177/life-tracker/backend/internal/repo/lifelog"
	"github.com/luyb177/life-tracker/backend/internal/repo/tag"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLifeLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListLifeLogLogic 生活记录列表（游标分页）
func NewListLifeLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLifeLogLogic {
	return &ListLifeLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLifeLogLogic) ListLifeLog(req *types.ListLifeLogReq) (resp *types.ListLifeLogResp, err error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	limit := int(req.PageSize)
	if limit <= 0 || limit > 50 {
		limit = constvar.DefaultPageSize
	}

	var cursorID uint64
	var cursorTime time.Time
	if req.PageToken != "" {
		var pt types.PageToken
		if err := pagetoken.Decode(req.PageToken, constvar.LifeLogPageTokenPrefix, &pt); err != nil {
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

	// 如果按标签过滤，先查 life_log IDs
	var filteredIDs []uint64
	if req.TagID > 0 {
		ids, err := l.svcCtx.Repos.Tag.FindLifeLogIDsByTagID(l.ctx, req.TagID, authUser.UserID)
		if err != nil {
			l.Errorf("find life log ids by tag failed: %v", err)
			return nil, errorx.WrapDBQuery("按标签查询失败", err)
		}
		filteredIDs = ids
	}

	// 多查一条判断 HasMore
	var logs []*lifelog.LifeLog
	if len(filteredIDs) > 0 || req.TagID > 0 {
		logs, err = l.svcCtx.Repos.LifeLog.ListByUserAndIDs(l.ctx, authUser.UserID, filteredIDs, cursorID, cursorTime, limit+1)
	} else {
		logs, err = l.svcCtx.Repos.LifeLog.ListByUser(l.ctx, authUser.UserID, cursorID, cursorTime, limit+1)
	}
	if err != nil {
		l.Errorf("list life logs failed: %v", err)
		return nil, errorx.WrapDBQuery("查询生活记录失败", err)
	}

	hasMore := len(logs) > limit
	if hasMore {
		logs = logs[:limit]
	}

	// 批量查标签
	lifeLogIDs := make([]uint64, 0, len(logs))
	for _, log := range logs {
		lifeLogIDs = append(lifeLogIDs, log.ID)
	}
	tagMap, err := l.svcCtx.Repos.Tag.BatchFindByLifeLogIDs(l.ctx, lifeLogIDs)
	if err != nil {
		l.Errorf("batch find tags failed: %v", err)
		return nil, errorx.WrapDBQuery("查询标签失败", err)
	}

	items := make([]types.LifeLogInfo, 0, len(logs))
	for _, log := range logs {
		tags := tagMap[log.ID]
		if tags == nil {
			tags = []*tag.Tag{}
		}
		tagInfos := make([]types.TagInfo, 0, len(tags))
		for _, t := range tags {
			tagInfos = append(tagInfos, types.TagInfo{ID: t.ID, Name: t.Name})
		}
		items = append(items, types.LifeLogInfo{
			ID:         log.ID,
			Content:    log.Content,
			Tags:       tagInfos,
			OccurredAt: log.OccurredAt.In(constvar.TimeLocation).Format(time.DateTime),
			CreatedAt:  log.CreatedAt.In(constvar.TimeLocation).Format(time.DateTime),
			UpdatedAt:  log.UpdatedAt.In(constvar.TimeLocation).Format(time.DateTime),
		})
	}

	var nextToken string
	if hasMore && len(logs) > 0 {
		last := logs[len(logs)-1]
		nextPT := types.PageToken{
			ID:        last.ID,
			CreatedAt: last.OccurredAt.In(constvar.TimeLocation).Format(time.DateTime),
		}
		nextToken, _ = pagetoken.Encode(constvar.LifeLogPageTokenPrefix, &nextPT)
	}

	return &types.ListLifeLogResp{
		List:      items,
		PageToken: nextToken,
		HasMore:   hasMore,
	}, nil
}
