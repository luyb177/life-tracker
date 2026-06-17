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

type ExpenseStatsRangeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExpenseStatsRangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExpenseStatsRangeLogic {
	return &ExpenseStatsRangeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExpenseStatsRangeLogic) ExpenseStatsRange(req *types.ExpenseStatsReq) (*types.ExpenseStatsRangeResp, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	start, end, err := parseDateRange(req.Start, req.End)
	if err != nil {
		return nil, err
	}

	total, err := l.svcCtx.Repos.Expense.SumByDateRange(l.ctx, authUser.UserID, start, end)
	if err != nil {
		l.Errorf("stats range failed: %v", err)
		return nil, errorx.WrapDBQuery("查询支出统计失败", err)
	}

	return &types.ExpenseStatsRangeResp{Total: total}, nil
}

func parseDateRange(startStr, endStr string) (time.Time, time.Time, error) {
	start, err := time.ParseInLocation("2006-01-02", startStr, constvar.TimeLocation)
	if err != nil {
		return time.Time{}, time.Time{}, errorx.WrapBadRequest("起始日期格式无效", err)
	}
	end, err := time.ParseInLocation("2006-01-02", endStr, constvar.TimeLocation)
	if err != nil {
		return time.Time{}, time.Time{}, errorx.WrapBadRequest("结束日期格式无效", err)
	}
	return start, end.Add(24 * time.Hour), nil
}
