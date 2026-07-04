// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package expense

import (
	"context"
	"time"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExpenseByDateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExpenseByDateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExpenseByDateLogic {
	return &ExpenseByDateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExpenseByDateLogic) ExpenseByDate(req *types.ExpenseByDateReq) (*types.ExpenseByDateResp, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	date, err := time.ParseInLocation("2006-01-02", req.Date, constvar.TimeLocation)
	if err != nil {
		return nil, errorx.WrapBadRequest("日期格式无效，期望: YYYY-MM-DD", err)
	}

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	logs, err := l.svcCtx.Repos.Expense.ListLogsByDateRange(l.ctx, authUser.UserID, startOfDay, endOfDay)
	if err != nil {
		l.Errorf("query expense by date failed: %v", err)
		return nil, errorx.WrapDBQuery("查询支出记录失败", err)
	}

	// 批量查询分类（含系统默认）
	categoryMap := make(map[uint64]string)
	categories, _ := l.svcCtx.Repos.Expense.FindCategoriesByUser(l.ctx, authUser.UserID)
	for _, c := range categories {
		categoryMap[c.ID] = c.Name
	}

	var total int64
	items := make([]types.ExpenseLogInfo, 0, len(logs))
	for _, log := range logs {
		total += log.Amount
		items = append(items, types.ExpenseLogInfo{
			ID: log.ID,
			Category: types.ExpenseCategoryInfo{
				ID:   log.CategoryID,
				Name: categoryMap[log.CategoryID],
			},
			Amount:     log.Amount,
			Note:       log.Note,
			Location:   log.Location,
			OccurredAt: log.OccurredAt.In(constvar.TimeLocation).Format(time.DateTime),
			Status:     log.Status,
			RefundedAt: formatTimePtr(log.RefundedAt, constvar.TimeLocation),
			CreatedAt:  log.CreatedAt.In(constvar.TimeLocation).Format(time.DateTime),
		})
	}

	return &types.ExpenseByDateResp{List: items, Total: total}, nil
}
