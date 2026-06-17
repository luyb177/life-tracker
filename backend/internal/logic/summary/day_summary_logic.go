// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package summary

import (
	"context"
	"time"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
	"github.com/luyb177/life-tracker/backend/internal/middleware"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DaySummaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDaySummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DaySummaryLogic {
	return &DaySummaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DaySummaryLogic) DaySummary(req *types.SummaryDayReq) (*types.SummaryDayResp, error) {
	authUser := middleware.GetAuthUser(l.ctx)
	if authUser == nil {
		return nil, errorx.ErrUnauthorized
	}

	if !reDay.MatchString(req.Date) {
		return nil, errorx.WrapBadRequest("日期格式无效，期望: YYYY-MM-DD", nil)
	}

	parsed, err := time.ParseInLocation("2006-01-02", req.Date, constvar.TimeLocation)
	if err != nil {
		return nil, errorx.WrapBadRequest("日期格式无效", err)
	}
	nextDate := parsed.AddDate(0, 0, 1)

	list, err := l.svcCtx.Repos.Summary.FindByPeriodRange(l.ctx, authUser.UserID, constvar.SummaryPeriodTypeDay, req.Date, nextDate.Format("2006-01-02"))
	if err != nil {
		l.Errorf("query day summary failed: %v", err)
		return nil, errorx.WrapDBQuery("查询日报失败", err)
	}

	items := make([]types.SummaryInfo, 0, len(list))
	for _, s := range list {
		items = append(items, types.SummaryInfo{
			ID:                s.ID,
			PeriodType:        s.PeriodType,
			PeriodStart:       s.PeriodStart,
			PeriodEnd:         s.PeriodEnd,
			Source:            s.Source,
			SummaryContent:    s.SummaryContent,
			SuggestionContent: s.SuggestionContent,
			Title:             s.Title,
			Tags:              s.Tags,
			Location:          s.Location,
			CreatedAt:         s.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:         s.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.SummaryDayResp{List: items}, nil
}
