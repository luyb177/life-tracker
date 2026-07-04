package lifelog

import (
	"context"
	"fmt"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"
)

// resolveTags 解析标签列表：id==0 则按名称创建标签，id>0 则验证标签存在
func resolveTags(ctx context.Context, svcCtx *svc.ServiceContext, tags []types.TagInfo) ([]uint64, error) {
	if len(tags) == 0 {
		return nil, nil
	}

	tagIDs := make([]uint64, 0, len(tags))
	for _, t := range tags {
		if t.ID == 0 {
			// 按名称查找或创建新标签
			tag, err := svcCtx.Repos.Tag.FindOrCreate(ctx, t.Name)
			if err != nil {
				return nil, errorx.WrapDBInsert(fmt.Sprintf("创建标签「%s」失败", t.Name), err)
			}
			tagIDs = append(tagIDs, tag.ID)
		} else {
			// 验证标签存在
			tag, err := svcCtx.Repos.Tag.FindByID(ctx, t.ID)
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
