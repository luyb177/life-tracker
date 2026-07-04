// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package lifelog

import (
	"context"
	"strings"
	"time"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLifeLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateLifeLogLogic 更新生活记录
func NewUpdateLifeLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLifeLogLogic {
	return &UpdateLifeLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLifeLogLogic) UpdateLifeLog(req *types.UpdateLifeLogReq) (resp *types.Response, err error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	existing, err := l.svcCtx.Repos.LifeLog.FindByID(l.ctx, req.ID)
	if err != nil {
		l.Errorf("find life log failed: %v", err)
		return nil, errorx.WrapDBQuery("查询生活记录失败", err)
	}
	if existing == nil {
		return nil, errorx.ErrNotFound
	}
	if existing.UserID != authUser.UserID {
		return nil, errorx.ErrForbidden
	}

	updates := make(map[string]interface{})
	if req.Content != "" {
		if len([]rune(req.Content)) > 10000 {
			return nil, errorx.WrapBadRequest("内容过长", nil)
		}
		updates["content"] = strings.TrimSpace(req.Content)
	}
	if req.OccurredAt != "" {
		occurredAt, err := time.ParseInLocation(time.DateTime, req.OccurredAt, constvar.TimeLocation)
		if err != nil {
			return nil, errorx.WrapBadRequest("时间格式无效", err)
		}
		updates["occurred_at"] = occurredAt
	}

	// 如果传了 tags，替换标签关联
	if req.Tags != nil {
		tagIDs, err := resolveTags(l.ctx, l.svcCtx, req.Tags)
		if err != nil {
			return nil, err
		}
		if err := l.svcCtx.Repos.Tag.DeleteByLifeLogID(l.ctx, req.ID); err != nil {
			l.Errorf("delete old tags failed: %v", err)
			return nil, errorx.WrapDBDelete("删除旧标签关联失败", err)
		}
		if len(tagIDs) > 0 {
			if err := l.svcCtx.Repos.Tag.BatchLink(l.ctx, req.ID, tagIDs); err != nil {
				l.Errorf("link tags failed: %v", err)
				return nil, errorx.WrapDBInsert("关联标签失败", err)
			}
		}
	}

	if len(updates) > 0 {
		if err := l.svcCtx.Repos.LifeLog.Update(l.ctx, req.ID, updates); err != nil {
			l.Errorf("update life log failed: %v", err)
			return nil, errorx.WrapDBUpdate("更新生活记录失败", err)
		}
	}

	return &types.Response{}, nil
}
