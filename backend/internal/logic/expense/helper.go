package expense

import (
	"time"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
	expenseRepo "github.com/luyb177/life-tracker/backend/internal/repo/expense"
	"github.com/luyb177/life-tracker/backend/internal/types"
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

func categoryInfoMap(categories []*expenseRepo.Category) map[uint64]types.ExpenseCategoryInfo {
	result := make(map[uint64]types.ExpenseCategoryInfo, len(categories))
	for _, c := range categories {
		result[c.ID] = types.ExpenseCategoryInfo{ID: c.ID, Name: c.Name, Type: c.Type}
	}
	return result
}

func expenseLogInfos(logs []*expenseRepo.Log, categoryMap map[uint64]types.ExpenseCategoryInfo) []types.ExpenseLogInfo {
	items := make([]types.ExpenseLogInfo, 0, len(logs))
	for _, log := range logs {
		items = append(items, expenseLogInfo(log, categoryMap[log.CategoryID]))
	}
	return items
}

func expenseLogInfo(log *expenseRepo.Log, category types.ExpenseCategoryInfo) types.ExpenseLogInfo {
	return types.ExpenseLogInfo{
		ID:            log.ID,
		Category:      category,
		Amount:        log.Amount,
		Note:          log.Note,
		Location:      log.Location,
		OccurredAt:    log.OccurredAt.In(constvar.TimeLocation).Format(time.DateTime),
		Status:        log.Status,
		RefundedAt:    formatTimePtr(log.RefundedAt, constvar.TimeLocation),
		CreatedAt:     log.CreatedAt.In(constvar.TimeLocation).Format(time.DateTime),
		UpdatedAt:     log.UpdatedAt.In(constvar.TimeLocation).Format(time.DateTime),
		LastUpdatedBy: log.LastUpdatedBy,
		LastUpdatedAt: formatTime(log.LastUpdatedAt, constvar.TimeLocation),
	}
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
