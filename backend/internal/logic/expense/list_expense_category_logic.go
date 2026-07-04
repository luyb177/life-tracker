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

type ListExpenseCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListExpenseCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListExpenseCategoryLogic {
	return &ListExpenseCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListExpenseCategoryLogic) ListExpenseCategory(req *types.Response) (*types.ListExpenseCategoryResp, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	categories, err := l.svcCtx.Repos.Expense.FindCategoriesByUser(l.ctx, authUser.UserID)
	if err != nil {
		l.Errorf("find categories failed: %v", err)
		return nil, errorx.WrapDBQuery("查询分类失败", err)
	}

	items := make([]types.ExpenseCategoryInfo, 0, len(categories))
	for _, c := range categories {
		items = append(items, types.ExpenseCategoryInfo{ID: c.ID, Name: c.Name, Type: c.Type})
	}

	return &types.ListExpenseCategoryResp{Categories: items}, nil
}
