package expense

import (
	"time"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
)

func formatTimePtr(t *time.Time, loc *time.Location) string {
	if t == nil {
		return ""
	}
	return t.In(loc).Format(time.DateTime)
}

func formatTime(t time.Time, loc *time.Location) string {
	if t.IsZero() {
		return ""
	}
	return t.In(loc).Format(time.DateTime)
}

func parseDateRange(startStr, endStr string) (time.Time, time.Time, error) {
	start, err := time.ParseInLocation("2006-01-02", startStr, constvar.TimeLocation)
	if err != nil {
		return time.Time{}, time.Time{}, errorx.WrapBadRequest("起始日期格式无效", err)
	}
	end, err := time.ParseInLocation("2006-01-02", endStr, constvar.TimeLocation)
	if err != nil {
		return time.Time{}, time.Time{}, errorx.WrapBadRequest("结束日期格式无效", err)
	}
	if end.Before(start) {
		return time.Time{}, time.Time{}, errorx.WrapBadRequest("结束日期不能早于起始日期", nil)
	}
	return start, end.Add(24 * time.Hour), nil
}
