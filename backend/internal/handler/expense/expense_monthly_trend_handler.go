// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package expense

import (
	"net/http"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/common/respx"
	"github.com/luyb177/life-tracker/backend/internal/logic/expense"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 每月支出趋势
func ExpenseMonthlyTrendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExpenseStatsReq
		if err := httpx.Parse(r, &req); err != nil {
			respx.ErrorCtx(r.Context(), w, errorx.WrapBadRequest("请求参数解析失败", err))
			return
		}

		l := expense.NewExpenseMonthlyTrendLogic(r.Context(), svcCtx)
		resp, err := l.ExpenseMonthlyTrend(&req)
		if err != nil {
			respx.ErrorCtx(r.Context(), w, err)
			return
		}
		respx.OkCtx(r.Context(), w, resp)
	}
}
