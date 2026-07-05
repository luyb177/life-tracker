package lifelog

import (
	"context"
	"fmt"
	"time"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
	lifelogRepo "github.com/luyb177/life-tracker/backend/internal/repo/lifelog"
	"github.com/luyb177/life-tracker/backend/internal/repo/tag"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"
	"gorm.io/gorm"
)

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(constvar.TimeLocation).Format(time.DateTime)
}

func lifeLogInfos(logs []*lifelogRepo.LifeLog, tagMap map[uint64][]*tag.Tag) []types.LifeLogInfo {
	items := make([]types.LifeLogInfo, 0, len(logs))
	for _, log := range logs {
		items = append(items, lifeLogInfo(log, tagMap[log.ID]))
	}
	return items
}

func lifeLogInfo(log *lifelogRepo.LifeLog, tags []*tag.Tag) types.LifeLogInfo {
	tagInfos := make([]types.TagInfo, 0, len(tags))
	for _, t := range tags {
		tagInfos = append(tagInfos, types.TagInfo{ID: t.ID, Name: t.Name})
	}
	return types.LifeLogInfo{
		ID:            log.ID,
		Content:       log.Content,
		Tags:          tagInfos,
		OccurredAt:    log.OccurredAt.In(constvar.TimeLocation).Format(time.DateTime),
		CreatedAt:     log.CreatedAt.In(constvar.TimeLocation).Format(time.DateTime),
		UpdatedAt:     log.UpdatedAt.In(constvar.TimeLocation).Format(time.DateTime),
		LastUpdatedBy: log.LastUpdatedBy,
		LastUpdatedAt: formatTime(log.LastUpdatedAt),
	}
}

// resolveTags 解析标签列表：id==0 则按名称创建标签，id>0 则验证标签存在
func resolveTags(ctx context.Context, svcCtx *svc.ServiceContext, tags []types.TagInfo, tx ...*gorm.DB) ([]uint64, error) {
	if len(tags) == 0 {
		return nil, nil
	}

	tagIDs := make([]uint64, 0, len(tags))
	for _, t := range tags {
		if t.ID == 0 {
			// 按名称查找或创建新标签
			tag, err := svcCtx.Repos.Tag.FindOrCreate(ctx, t.Name, tx...)
			if err != nil {
				return nil, errorx.WrapDBInsert(fmt.Sprintf("创建标签「%s」失败", t.Name), err)
			}
			tagIDs = append(tagIDs, tag.ID)
		} else {
			// 验证标签存在
			tag, err := svcCtx.Repos.Tag.FindByID(ctx, t.ID, tx...)
			if err != nil {
				return nil, errorx.WrapDBQuery(fmt.Sprintf("查询标签失败"), err)
			}
			if tag == nil {
				return nil, errorx.ErrNotFound
			}
			tagIDs = append(tagIDs, tag.ID)
		}
	}
	return tagIDs, nil
}
