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

type ExpenseStatsCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExpenseStatsCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExpenseStatsCategoryLogic {
	return &ExpenseStatsCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExpenseStatsCategoryLogic) ExpenseStatsCategory(req *types.ExpenseStatsReq) (*types.ExpenseStatsCategoryResp, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	start, end, err := parseDateRange(req.Start, req.End)
	if err != nil {
		return nil, err
	}

	breakdown, err := l.svcCtx.Repos.Expense.SumByDateRangeGrouped(l.ctx, authUser.UserID, start, end)
	if err != nil {
		l.Errorf("stats category failed: %v", err)
		return nil, errorx.WrapDBQuery("查询分类统计失败", err)
	}

	categories := make([]types.ExpenseCategoryStat, 0, len(breakdown))
	for _, b := range breakdown {
		categories = append(categories, types.ExpenseCategoryStat{
			CategoryID:   b.CategoryID,
			CategoryName: b.CategoryName,
			Total:        b.Total,
		})
	}

	return &types.ExpenseStatsCategoryResp{Categories: categories}, nil
}
