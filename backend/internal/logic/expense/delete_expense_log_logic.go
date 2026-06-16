// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package expense

import (
	"context"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteExpenseLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteExpenseLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteExpenseLogLogic {
	return &DeleteExpenseLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteExpenseLogLogic) DeleteExpenseLog(req *types.DeleteExpenseLogReq) (*types.Response, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	log, err := l.svcCtx.Repos.Expense.FindLogByID(l.ctx, req.ID)
	if err != nil {
		l.Errorf("find expense log failed: %v", err)
		return nil, errorx.WrapDBQuery("查询记录失败", err)
	}
	if log == nil {
		return nil, errorx.ErrNotFound
	}
	if log.UserID != authUser.UserID {
		return nil, errorx.ErrForbidden
	}

	if err := l.svcCtx.Repos.Expense.DeleteLog(l.ctx, req.ID); err != nil {
		l.Errorf("delete expense log failed: %v", err)
		return nil, errorx.WrapDBDelete("删除记录失败", err)
	}

	return &types.Response{}, nil
}
