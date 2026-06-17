// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package expense

import (
	"context"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteExpenseCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteExpenseCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteExpenseCategoryLogic {
	return &DeleteExpenseCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteExpenseCategoryLogic) DeleteExpenseCategory(req *types.DeleteExpenseCategoryReq) (*types.Response, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	c, err := l.svcCtx.Repos.Expense.FindCategoryByID(l.ctx, req.ID)
	if err != nil {
		l.Errorf("find category failed: %v", err)
		return nil, errorx.WrapDBQuery("查询分类失败", err)
	}
	if c == nil {
		return nil, errorx.ErrNotFound
	}
	if c.UserID != authUser.UserID {
		return nil, errorx.ErrForbidden
	}
	// 系统默认分类不可删除
	if c.Type == constvar.ExpenseCategoryTypeSystem {
		return nil, errorx.WrapBadRequest("系统默认分类不可删除", nil)
	}

	// 检查是否存在关联的支出记录
	count, err := l.svcCtx.Repos.Expense.CountLogsByCategory(l.ctx, req.ID)
	if err != nil {
		l.Errorf("count logs by category failed: %v", err)
		return nil, errorx.WrapDBQuery("查询关联记录失败", err)
	}
	if count > 0 {
		return nil, errorx.WrapBadRequest("该分类下存在支出记录，无法删除", nil)
	}

	if err := l.svcCtx.Repos.Expense.DeleteCategory(l.ctx, req.ID); err != nil {
		l.Errorf("delete category failed: %v", err)
		return nil, errorx.WrapDBDelete("删除分类失败", err)
	}

	return &types.Response{}, nil
}
