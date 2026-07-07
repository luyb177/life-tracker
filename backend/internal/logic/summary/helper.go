package summary

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/luyb177/life-tracker/backend/common/errorx"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
	"github.com/luyb177/life-tracker/backend/internal/repo/summary"
	"github.com/luyb177/life-tracker/backend/internal/svc"
	"github.com/luyb177/life-tracker/backend/internal/types"
	"gorm.io/gorm"
)

var reDay = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

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

	if !validPeriodType(periodType) {
		return time.Time{}, fmt.Errorf("无效的周期类型")
	}
	return startAt, nil
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
	default:
		return "总结"
	}
}

func periodStartHint(t uint8) string {
	return "YYYY-MM-DD"
}

// resolveSummaryTags 解析标签列表并关联到 summary
func resolveSummaryTags(ctx context.Context, svcCtx *svc.ServiceContext, summaryID uint64, tags []types.TagInfo, tx ...*gorm.DB) error {
	if len(tags) == 0 {
		return nil
	}
	tagIDs := make([]uint64, 0, len(tags))
	for _, t := range tags {
		if t.ID == 0 {
			tag, err := svcCtx.Repos.Tag.FindOrCreate(ctx, t.Name, tx...)
			if err != nil {
				return errorx.WrapDBInsert("创建标签失败", err)
			}
			tagIDs = append(tagIDs, tag.ID)
		} else {
			tag, err := svcCtx.Repos.Tag.FindByID(ctx, t.ID, tx...)
			if err != nil {
				return errorx.WrapDBQuery("查询标签失败", err)
			}
			if tag == nil {
				return errorx.ErrNotFound
			}
			tagIDs = append(tagIDs, tag.ID)
		}
	}
	return svcCtx.Repos.Tag.BatchLinkSummary(ctx, summaryID, tagIDs, tx...)
}

// buildExpenseLocationSummary keeps manual summaries aligned with AI summaries:
// the location field describes where spending happened during the period.
func buildExpenseLocationSummary(ctx context.Context, svcCtx *svc.ServiceContext, userID uint64, start, end time.Time) string {
	logs, err := svcCtx.Repos.Expense.ListLogsByDateRange(ctx, userID, start, end)
	if err != nil || len(logs) == 0 {
		return ""
	}

	locTotals := make(map[string]int64)
	for _, log := range logs {
		if log.Status == 1 || log.Amount <= 0 {
			continue
		}
		loc := strings.TrimSpace(log.Location)
		if loc == "" || loc == "未知" {
			continue
		}
		locTotals[loc] += log.Amount
	}
	if len(locTotals) == 0 {
		return ""
	}

	locations := make([]string, 0, len(locTotals))
	for loc := range locTotals {
		locations = append(locations, loc)
	}
	sort.Strings(locations)

	parts := make([]string, 0, len(locations))
	for _, loc := range locations {
		parts = append(parts, fmt.Sprintf("%s：%.2f 元", loc, float64(locTotals[loc])/100))
	}
	return strings.Join(parts, "；")
}

// batchFillSummaryTags 批量填充 summary 的标签
func batchFillSummaryTags(ctx context.Context, svcCtx *svc.ServiceContext, summaryIDs []uint64) (map[uint64][]types.TagInfo, error) {
	tagMap, err := svcCtx.Repos.Tag.BatchFindBySummaryIDs(ctx, summaryIDs)
	if err != nil {
		return nil, err
	}
	result := make(map[uint64][]types.TagInfo)
	for id, tags := range tagMap {
		infos := make([]types.TagInfo, 0, len(tags))
		for _, t := range tags {
			infos = append(infos, types.TagInfo{ID: t.ID, Name: t.Name})
		}
		result[id] = infos
	}
	return result, nil
}

// summaryToInfo 将 summary 模型转为 API 响应
func summaryToInfo(s *summary.Summary, tagInfos []types.TagInfo) types.SummaryInfo {
	return types.SummaryInfo{
		ID:                s.ID,
		PeriodType:        s.PeriodType,
		PeriodStart:       s.PeriodStart,
		PeriodEnd:         s.PeriodEnd,
		Source:            s.Source,
		SummaryContent:    s.SummaryContent,
		SuggestionContent: s.SuggestionContent,
		Title:             s.Title,
		Tags:              tagInfos,
		Location:          s.Location,
		CreatedAt:         s.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         s.UpdatedAt.Format("2006-01-02 15:04:05"),
		LastUpdatedBy:     s.LastUpdatedBy,
		LastUpdatedAt:     formatSummaryTime(s.LastUpdatedAt),
	}
}

func formatSummaryTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(constvar.TimeLocation).Format(time.DateTime)
}
