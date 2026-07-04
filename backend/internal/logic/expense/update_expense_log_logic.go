// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package expense

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

type UpdateExpenseLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateExpenseLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateExpenseLogLogic {
	return &UpdateExpenseLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateExpenseLogLogic) UpdateExpenseLog(req *types.UpdateExpenseLogReq) (*types.Response, error) {
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

	updates := make(map[string]interface{})
	if req.CategoryID > 0 {
		cat, err := l.svcCtx.Repos.Expense.FindCategoryByID(l.ctx, req.CategoryID)
		if err != nil {
			return nil, errorx.WrapDBQuery("查询分类失败", err)
		}
		if cat == nil || (cat.UserID != 0 && cat.UserID != authUser.UserID) {
			return nil, errorx.WrapBadRequest("分类无效", nil)
		}
		updates["category_id"] = req.CategoryID
	}
	if req.Amount > 0 {
		updates["amount"] = req.Amount
	}
	if req.Note != nil {
		if len([]rune(*req.Note)) > 255 {
			return nil, errorx.WrapBadRequest("备注过长", nil)
		}
		updates["note"] = strings.TrimSpace(*req.Note)
	}
	if req.OccurredAt != "" {
		occurredAt, err := time.ParseInLocation(time.DateTime, req.OccurredAt, constvar.TimeLocation)
		if err != nil {
			return nil, errorx.WrapBadRequest("时间格式无效", err)
		}
		updates["occurred_at"] = occurredAt
	}
	if len(updates) == 0 {
		return nil, errorx.WrapBadRequest("没有可更新的字段", nil)
	}

	if err := l.svcCtx.Repos.Expense.UpdateLog(l.ctx, req.ID, updates); err != nil {
		l.Errorf("update expense log failed: %v", err)
		return nil, errorx.WrapDBUpdate("更新支出记录失败", err)
	}

	return &types.Response{}, nil
}
