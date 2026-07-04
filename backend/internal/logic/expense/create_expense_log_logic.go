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
	"github.com/luyb177/life-tracker/backend/internal/repo/expense"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateExpenseLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateExpenseLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateExpenseLogLogic {
	return &CreateExpenseLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateExpenseLogLogic) CreateExpenseLog(req *types.CreateExpenseLogReq) (*types.Response, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	if req.Amount <= 0 {
		return nil, errorx.WrapBadRequest("金额必须大于0（单位：分）", nil)
	}
	if req.CategoryID == 0 {
		return nil, errorx.WrapBadRequest("请选择分类", nil)
	}

	category, err := l.svcCtx.Repos.Expense.FindCategoryByID(l.ctx, req.CategoryID)
	if err != nil {
		l.Errorf("find expense category failed: %v", err)
		return nil, errorx.WrapDBQuery("查询分类失败", err)
	}
	if category == nil {
		return nil, errorx.ErrNotFound
	}
	// 系统默认(user_id=0)或用户本人的分类才允许
	if category.UserID != 0 && category.UserID != authUser.UserID {
		return nil, errorx.ErrForbidden
	}

	if req.Note != nil && len([]rune(*req.Note)) > 255 {
		return nil, errorx.WrapBadRequest("备注过长", nil)
	}

	occurredAt, err := time.ParseInLocation(time.DateTime, req.OccurredAt, constvar.TimeLocation)
	if err != nil {
		return nil, errorx.WrapBadRequest("时间格式无效", err)
	}

	// 从 IP 中间件获取地理位置
	locStr := middleware.FullLocation(middleware.GetIPLocation(l.ctx))

	note := ""
	if req.Note != nil {
		note = strings.TrimSpace(*req.Note)
	}

	log := &expense.Log{
		UserID:     authUser.UserID,
		CategoryID: req.CategoryID,
		Amount:     req.Amount,
		Note:       note,
		Location:   locStr,
		OccurredAt: occurredAt,
	}

	if err := l.svcCtx.Repos.Expense.CreateLog(l.ctx, log); err != nil {
		l.Errorf("create expense log failed: %v", err)
		return nil, errorx.WrapDBInsert("记录支出失败", err)
	}

	return &types.Response{}, nil
}
