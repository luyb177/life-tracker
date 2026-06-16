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

type DailyExpenseTotalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDailyExpenseTotalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DailyExpenseTotalLogic {
	return &DailyExpenseTotalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DailyExpenseTotalLogic) DailyExpenseTotal(req *types.DailyExpenseTotalReq) (*types.DailyExpenseTotalResp, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	date, err := time.ParseInLocation(time.DateOnly, req.Date, constvar.TimeLocation)
	if err != nil {
		return nil, errorx.WrapBadRequest("日期格式无效", err)
	}

	total, err := l.svcCtx.Repos.Expense.SumByDate(l.ctx, authUser.UserID, date)
	if err != nil {
		l.Errorf("sum expense failed: %v", err)
		return nil, errorx.WrapDBQuery("汇总支出失败", err)
	}

	return &types.DailyExpenseTotalResp{
		Date:  req.Date,
		Total: total,
	}, nil
}
