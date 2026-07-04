// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package lifelog

import (
	"net/http"

	"github.com/luyb177/life-tracker/backend/common/respx"
	"github.com/luyb177/life-tracker/backend/internal/logic/lifelog"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 生活记录列表（游标分页）
func ListLifeLogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListLifeLogReq
		if err := httpx.Parse(r, &req); err != nil {
			respx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := lifelog.NewListLifeLogLogic(r.Context(), svcCtx)
		resp, err := l.ListLifeLog(&req)
		if err != nil {
			respx.ErrorCtx(r.Context(), w, err)
		} else {
			respx.OkCtx(r.Context(), w, resp)
		}
	}
}
