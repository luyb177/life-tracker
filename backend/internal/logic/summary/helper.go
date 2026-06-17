package summary

import (
	"fmt"
	"regexp"
	"time"

	"github.com/luyb177/life-tracker/backend/internal/constvar"
)

var reDay = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func validPeriodType(t uint8) bool {
	switch t {
	case constvar.SummaryPeriodTypeDay,
		constvar.SummaryPeriodTypeWeek,
		constvar.SummaryPeriodTypeMonth,
		constvar.SummaryPeriodTypeYear,
		constvar.SummaryPeriodTypeLife:
		return true
	default:
		return false
	}
}

func parsePeriodDate(value string) (time.Time, error) {
	if !reDay.MatchString(value) {
		return time.Time{}, fmt.Errorf("日期格式无效，期望: YYYY-MM-DD")
	}
	parsed, err := time.ParseInLocation(time.DateOnly, value, constvar.TimeLocation)
	if err != nil {
		return time.Time{}, fmt.Errorf("日期格式无效: %w", err)
	}
	return time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, constvar.TimeLocation), nil
}

func normalizePeriodStart(periodType uint8, start string) (time.Time, error) {
	startAt, err := parsePeriodDate(start)
	if err != nil {
		return time.Time{}, err
	}

	switch periodType {
	case constvar.SummaryPeriodTypeDay:
		return startAt, nil
	case constvar.SummaryPeriodTypeWeek:
		if startAt.Weekday() != time.Monday {
			return time.Time{}, fmt.Errorf("周报的 period_start 必须是周一")
		}
		return startAt, nil
	case constvar.SummaryPeriodTypeMonth:
		if startAt.Day() != 1 {
			return time.Time{}, fmt.Errorf("月报的 period_start 必须是每月第一天")
		}
		return startAt, nil
	case constvar.SummaryPeriodTypeYear:
		if startAt.Month() != time.January || startAt.Day() != 1 {
			return time.Time{}, fmt.Errorf("年报的 period_start 必须是每年 1 月 1 日")
		}
		return startAt, nil
	case constvar.SummaryPeriodTypeLife:
		return startAt, nil
	default:
		return time.Time{}, fmt.Errorf("无效的周期类型")
	}
}

func normalizePeriodRange(periodType uint8, start, end string) (time.Time, time.Time, error) {
	startAt, err := normalizePeriodStart(periodType, start)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	endAt, err := parsePeriodDate(end)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	if !endAt.After(startAt) {
		return time.Time{}, time.Time{}, fmt.Errorf("period_end 必须晚于 period_start")
	}

	switch periodType {
	case constvar.SummaryPeriodTypeDay:
		if !endAt.Equal(startAt.AddDate(0, 0, 1)) {
			return time.Time{}, time.Time{}, fmt.Errorf("日报必须满足 period_end = period_start + 1 天")
		}
	case constvar.SummaryPeriodTypeWeek:
		if !endAt.Equal(startAt.AddDate(0, 0, 7)) {
			return time.Time{}, time.Time{}, fmt.Errorf("周报必须满足 period_end = period_start + 7 天")
		}
	case constvar.SummaryPeriodTypeMonth:
		if !endAt.Equal(startAt.AddDate(0, 1, 0)) {
			return time.Time{}, time.Time{}, fmt.Errorf("月报必须满足 period_end = 下月第一天")
		}
	case constvar.SummaryPeriodTypeYear:
		if !endAt.Equal(startAt.AddDate(1, 0, 0)) {
			return time.Time{}, time.Time{}, fmt.Errorf("年报必须满足 period_end = 下一年 1 月 1 日")
		}
	case constvar.SummaryPeriodTypeLife:
		// 人生总结允许用户自定义时间范围，只要求 end > start。
	default:
		return time.Time{}, time.Time{}, fmt.Errorf("无效的周期类型")
	}

	return startAt, endAt, nil
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
	case constvar.SummaryPeriodTypeLife:
		return "人生总结"
	default:
		return "总结"
	}
}

func periodStartHint(t uint8) string {
	return "YYYY-MM-DD"
}
