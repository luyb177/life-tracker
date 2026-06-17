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

type ExpenseStatsTrendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExpenseStatsTrendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExpenseStatsTrendLogic {
	return &ExpenseStatsTrendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExpenseStatsTrendLogic) ExpenseStatsTrend(req *types.ExpenseStatsReq) (*types.ExpenseStatsTrendResp, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	start, end, err := parseDateRange(req.Start, req.End)
	if err != nil {
		return nil, err
	}

	days, err := l.svcCtx.Repos.Expense.SumByDay(l.ctx, authUser.UserID, start, end)
	if err != nil {
		l.Errorf("stats trend failed: %v", err)
		return nil, errorx.WrapDBQuery("查询趋势数据失败", err)
	}

	points := make([]types.ExpenseTrendPoint, 0, len(days))
	for _, d := range days {
		points = append(points, types.ExpenseTrendPoint{
			Date:  d.Date,
			Total: d.Total,
		})
	}

	return &types.ExpenseStatsTrendResp{Points: points}, nil
}
