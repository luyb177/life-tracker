// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package expense

import (
	"context"
	"time"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/pkg/pagetoken"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListExpenseLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListExpenseLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListExpenseLogLogic {
	return &ListExpenseLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListExpenseLogLogic) ListExpenseLog(req *types.ListExpenseLogReq) (*types.ListExpenseLogResp, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	limit := int(req.PageSize)
	if limit <= 0 || limit > 50 {
		limit = constvar.DefaultPageSize
	}

	var cursorID uint64
	var cursorTime time.Time
	if req.PageToken != "" {
		var pt types.PageToken
		if err := pagetoken.Decode(req.PageToken, constvar.ExpensePageTokenPrefix, &pt); err != nil {
			return nil, errorx.WrapBadRequest("分页参数无效", err)
		}
		cursorID = pt.ID
		if pt.CreatedAt != "" {
			var parseErr error
			cursorTime, parseErr = time.ParseInLocation(time.DateTime, pt.CreatedAt, constvar.TimeLocation)
			if parseErr != nil {
				return nil, errorx.WrapBadRequest("分页参数无效", parseErr)
			}
		}
	}

	// 多查一条判断 HasMore
	logs, err := l.svcCtx.Repos.Expense.ListLogsByUser(l.ctx, authUser.UserID, cursorID, cursorTime, limit+1)
	if err != nil {
		l.Errorf("list expense logs failed: %v", err)
		return nil, errorx.WrapDBQuery("查询支出记录失败", err)
	}

	hasMore := len(logs) > limit
	if hasMore {
		logs = logs[:limit]
	}

	// 批量查询分类
	categoryIDs := make([]uint64, 0, len(logs))
	for _, log := range logs {
		categoryIDs = append(categoryIDs, log.CategoryID)
	}
	categoryMap := make(map[uint64]string)
	if len(categoryIDs) > 0 {
		categories, err := l.svcCtx.Repos.Expense.FindCategoriesByUser(l.ctx, authUser.UserID)
		if err != nil {
			l.Errorf("find expense categories failed: %v", err)
			return nil, errorx.WrapDBQuery("查询分类失败", err)
		}
		for _, c := range categories {
			categoryMap[c.ID] = c.Name
		}
	}

	items := make([]types.ExpenseLogInfo, 0, len(logs))
	for _, log := range logs {
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
			CreatedAt:  log.CreatedAt.In(constvar.TimeLocation).Format(time.DateTime),
		})
	}

	var nextToken string
	if hasMore && len(logs) > 0 {
		last := logs[len(logs)-1]
		nextPT := types.PageToken{
			ID:        last.ID,
			CreatedAt: last.OccurredAt.In(constvar.TimeLocation).Format(time.DateTime),
		}
		nextToken, _ = pagetoken.Encode(constvar.ExpensePageTokenPrefix, &nextPT)
	}

	return &types.ListExpenseLogResp{
		List:      items,
		PageToken: nextToken,
		HasMore:   hasMore,
	}, nil
}
