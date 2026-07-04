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

// 按天查询生活记录
func LifeLogByDateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LifeLogByDateReq
		if err := httpx.Parse(r, &req); err != nil {
			respx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := lifelog.NewLifeLogByDateLogic(r.Context(), svcCtx)
		resp, err := l.LifeLogByDate(&req)
		if err != nil {
			respx.ErrorCtx(r.Context(), w, err)
		} else {
			respx.OkCtx(r.Context(), w, resp)
		}
	}
}
