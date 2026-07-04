package expense

import (
	"context"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundExpenseLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundExpenseLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundExpenseLogLogic {
	return &RefundExpenseLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefundExpenseLogLogic) RefundExpenseLog(req *types.RefundExpenseLogReq) (*types.Response, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	log, err := l.svcCtx.Repos.Expense.FindLogByID(l.ctx, req.ID)
	if err != nil {
		l.Errorf("find expense log failed: %v", err)
		return nil, errorx.WrapDBQuery("查询支出记录失败", err)
	}
	if log == nil {
		return nil, errorx.ErrNotFound
	}
	if log.UserID != authUser.UserID {
		return nil, errorx.ErrForbidden
	}
	if log.Status == 1 {
		return nil, errorx.WrapBadRequest("该记录已退款", nil)
	}

	if err := l.svcCtx.Repos.Expense.RefundLog(l.ctx, req.ID); err != nil {
		l.Errorf("refund expense log failed: %v", err)
		return nil, errorx.WrapDBUpdate("退款失败", err)
	}

	return &types.Response{}, nil
}
