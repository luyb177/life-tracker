// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package summary

import (
	"net/http"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/common/respx"
	"github.com/luyb177/life-tracker/backend/internal/logic/summary"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 查询某天日报（用户+AI）
func DaySummaryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SummaryDayReq
		if err := httpx.Parse(r, &req); err != nil {
			respx.ErrorCtx(r.Context(), w, errorx.WrapBadRequest("请求参数解析失败", err))
			return
		}

		l := summary.NewDaySummaryLogic(r.Context(), svcCtx)
		resp, err := l.DaySummary(&req)
		if err != nil {
			respx.ErrorCtx(r.Context(), w, err)
		} else {
			respx.OkCtx(r.Context(), w, resp)
		}
	}
}
