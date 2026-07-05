// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package lifelog

import (
	"context"
	"time"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/repo/tag"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LifeLogByDateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewLifeLogByDateLogic 按天查询生活记录
func NewLifeLogByDateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LifeLogByDateLogic {
	return &LifeLogByDateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LifeLogByDateLogic) LifeLogByDate(req *types.LifeLogByDateReq) (resp *types.LifeLogByDateResp, err error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	date, err := time.ParseInLocation("2006-01-02", req.Date, constvar.TimeLocation)
	if err != nil {
		return nil, errorx.WrapBadRequest("日期格式无效，期望: YYYY-MM-DD", err)
	}
	nextDate := date.AddDate(0, 0, 1)

	logs, err := l.svcCtx.Repos.LifeLog.FindByDateRange(l.ctx, authUser.UserID, date, nextDate)
	if err != nil {
		l.Errorf("find life logs by date failed: %v", err)
		return nil, errorx.WrapDBQuery("查询生活记录失败", err)
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
			ID:            log.ID,
			Content:       log.Content,
			Tags:          tagInfos,
			OccurredAt:    log.OccurredAt.In(constvar.TimeLocation).Format(time.DateTime),
			CreatedAt:     log.CreatedAt.In(constvar.TimeLocation).Format(time.DateTime),
			UpdatedAt:     log.UpdatedAt.In(constvar.TimeLocation).Format(time.DateTime),
			LastUpdatedBy: log.LastUpdatedBy,
			LastUpdatedAt: formatTime(log.LastUpdatedAt),
		})
	}

	return &types.LifeLogByDateResp{
		List: items,
	}, nil
}
