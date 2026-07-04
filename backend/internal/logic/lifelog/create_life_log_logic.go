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
	"github.com/luyb177/life-tracker/backend/internal/repo/lifelog"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLifeLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建生活记录
func NewCreateLifeLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLifeLogLogic {
	return &CreateLifeLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLifeLogLogic) CreateLifeLog(req *types.CreateLifeLogReq) (resp *types.IDResponse, err error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	if strings.TrimSpace(req.Content) == "" {
		return nil, errorx.WrapBadRequest("内容不能为空", nil)
	}
	if len([]rune(req.Content)) > 10000 {
		return nil, errorx.WrapBadRequest("内容过长", nil)
	}

	occurredAt, err := time.ParseInLocation(time.DateTime, req.OccurredAt, constvar.TimeLocation)
	if err != nil {
		return nil, errorx.WrapBadRequest("时间格式无效", err)
	}

	log := &lifelog.LifeLog{
		UserID:     authUser.UserID,
		Content:    strings.TrimSpace(req.Content),
		Tags:       strings.TrimSpace(req.Tags),
		OccurredAt: occurredAt,
	}

	if err := l.svcCtx.Repos.LifeLog.Create(l.ctx, log); err != nil {
		l.Errorf("create life log failed: %v", err)
		return nil, errorx.WrapDBInsert("创建生活记录失败", err)
	}

	return &types.IDResponse{ID: log.ID}, nil
}
