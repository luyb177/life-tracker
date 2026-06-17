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

// 按时间范围查询总结
func RangeSummaryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SummaryRangeReq
		if err := httpx.Parse(r, &req); err != nil {
			respx.ErrorCtx(r.Context(), w, errorx.WrapBadRequest("请求参数解析失败", err))
			return
		}

		l := summary.NewRangeSummaryLogic(r.Context(), svcCtx)
		resp, err := l.RangeSummary(&req)
		if err != nil {
			respx.ErrorCtx(r.Context(), w, err)
		} else {
			respx.OkCtx(r.Context(), w, resp)
		}
	}
}
