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

type ExpenseMonthlyTrendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExpenseMonthlyTrendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExpenseMonthlyTrendLogic {
	return &ExpenseMonthlyTrendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExpenseMonthlyTrendLogic) ExpenseMonthlyTrend(req *types.ExpenseStatsReq) (*types.ExpenseMonthlyTrendResp, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	start, end, err := parseDateRange(req.Start, req.End)
	if err != nil {
		return nil, err
	}

	months, err := l.svcCtx.Repos.Expense.SumByMonth(l.ctx, authUser.UserID, start, end)
	if err != nil {
		l.Errorf("monthly trend failed: %v", err)
		return nil, errorx.WrapDBQuery("查询月度趋势失败", err)
	}

	points := make([]types.MonthTotal, 0, len(months))
	for _, m := range months {
		points = append(points, types.MonthTotal{
			Month: m.Month,
			Total: m.Total,
		})
	}

	return &types.ExpenseMonthlyTrendResp{Points: points}, nil
}
