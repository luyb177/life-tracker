// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package expense

import (
	"context"
	"strings"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/repo/expense"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateExpenseCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateExpenseCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateExpenseCategoryLogic {
	return &CreateExpenseCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateExpenseCategoryLogic) CreateExpenseCategory(req *types.CreateExpenseCategoryReq) (*types.Response, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	if strings.TrimSpace(req.Name) == "" {
		return nil, errorx.WrapBadRequest("分类名称不能为空", nil)
	}

	c := &expense.Category{
		UserID: authUser.UserID,
		Name:   strings.TrimSpace(req.Name),
		Type:   constvar.ExpenseCategoryTypeUser,
	}

	if err := l.svcCtx.Repos.Expense.CreateCategory(l.ctx, c); err != nil {
		l.Errorf("create category failed: %v", err)
		return nil, errorx.WrapDBInsert("创建分类失败", err)
	}

	return &types.Response{}, nil
}
