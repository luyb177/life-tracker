package summary

import "github.com/luyb177/life-tracker/backend/internal/constvar"

func validPeriodType(t uint8) bool {
	switch t {
	case constvar.SummaryPeriodTypeDay,
		constvar.SummaryPeriodTypeWeek,
		constvar.SummaryPeriodTypeMonth,
		constvar.SummaryPeriodTypeYear:
		return true
	default:
		return false
	}
}

func periodTypeLabel(t uint8) string {
	switch t {
	case constvar.SummaryPeriodTypeDay:
		return "日报"
	case constvar.SummaryPeriodTypeWeek:
		return "周报"
	case constvar.SummaryPeriodTypeMonth:
		return "月报"
	case constvar.SummaryPeriodTypeYear:
		return "年报"
	default:
		return "总结"
	}
}
