package summary

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/luyb177/life-tracker/backend/internal/constvar"
	"github.com/luyb177/life-tracker/backend/internal/pkg/ai"
	expenseRepo "github.com/luyb177/life-tracker/backend/internal/repo/expense"
	"github.com/luyb177/life-tracker/backend/internal/repo/summary"
	"github.com/luyb177/life-tracker/backend/internal/svc"
)

// Run 执行 AI 总结
// periodStart 为零值时自动按当前时间计算周期；非零值时以此日期为起始计算结束日期。
func Run(ctx context.Context, svcCtx *svc.ServiceContext, periodType uint8, userID uint64, periodStart time.Time) error {
	var periodEnd time.Time
	if periodStart.IsZero() {
		periodStart, periodEnd = calcPeriod(periodType)
	} else {
		periodEnd = calcPeriodEnd(periodType, periodStart)
	}

	// 1. 查分类支出
	categoryBreakdown, err := svcCtx.Repos.Expense.SumByDateRangeGrouped(ctx, userID, periodStart, periodEnd)
	if err != nil {
		return fmt.Errorf("query expense grouped: %w", err)
	}
	totalExpense, err := svcCtx.Repos.Expense.SumByDateRange(ctx, userID, periodStart, periodEnd)
	if err != nil {
		return fmt.Errorf("query expense total: %w", err)
	}

	// 1.5 查地点分布
	locationBreakdown := buildLocationBreakdown(ctx, svcCtx, userID, periodStart, periodEnd)

	// 2. 构建上下文（下级总结）
	contextText := buildContext(ctx, svcCtx, periodType, userID, periodStart, periodEnd)

	// 2.5 获取用户日报（生活记录）
	journalText := buildJournalContext(ctx, svcCtx, userID, periodStart, periodEnd)

	// 2.6 人生总结额外上下文（标签统计 + 年报聚合）
	extraContext := ""
	if periodType == constvar.SummaryPeriodTypeLife {
		extraContext = buildLifetimeContext(ctx, svcCtx, userID, periodStart, periodEnd)
	}

	// 3. 构建 Prompt
	label := periodTypeLabelCN(periodType)
	systemPrompt := "你是一个个人生活助手。请根据提供的数据生成简洁的「" + label + "」总结。用中文回复，语气亲切。"
	userPrompt := fmt.Sprintf(`请生成%s总结：

周期：%s ~ %s

【支出汇总】
总支出：%.2f 元
分类明细：
%s
【地点分布】
%s
【用户记录】
%s
%s
【上下文】
%s

请包含：1. 消费概况与分类分析 2. 结合用户记录分析当日行为 3. 地点分布分析 4. 变化趋势（对比上下文） 5. 改进建议。
如果数据为空或极少，请如实说明"今日无事"或"本周期无记录"，不用强行编造。`,
		label,
		periodStart.Format("2006-01-02"), periodEnd.Format("2006-01-02"),
		totalExpense,
		formatCategoryBreakdown(categoryBreakdown),
		locationBreakdown,
		journalText,
		extraContext,
		contextText,
	)

	// 4. 调用 DeepSeek
	aiClient := ai.NewClient(svcCtx.Config.AIConf.Endpoint, svcCtx.Config.AIConf.APIKey, svcCtx.Config.AIConf.Model)
	aiContent, err := aiClient.Chat(ctx, systemPrompt, userPrompt)
	if err != nil {
		return fmt.Errorf("ai chat: %w", err)
	}

	// 5. 保存（按 source=AI 查找避免用户手写记录干扰去重）
	periodStartStr := periodStart.Format("2006-01-02")
	s := &summary.Summary{
		UserID:            userID,
		PeriodType:        periodType,
		PeriodStart:       periodStartStr,
		PeriodEnd:         periodEnd.Format("2006-01-02"),
		Source:            constvar.SummarySourceAI,
		SummaryContent:    aiContent,
		SuggestionContent: "",
		Location:          locationBreakdown,
	}

	existingAI, err := svcCtx.Repos.Summary.FindByPeriodAndSource(ctx, userID, periodType, periodStartStr, constvar.SummarySourceAI)
	if err != nil {
		return fmt.Errorf("find existing ai summary: %w", err)
	}
	if existingAI != nil {
		return svcCtx.Repos.Summary.Update(ctx, existingAI.ID, map[string]interface{}{
			"summary_content":    aiContent,
			"suggestion_content": "",
			"location":           locationBreakdown,
		})
	}
	return svcCtx.Repos.Summary.Create(ctx, s)
}

// buildJournalContext 获取周期内的用户日报内容作为 AI 上下文
func buildJournalContext(ctx context.Context, svcCtx *svc.ServiceContext, userID uint64, start, end time.Time) string {
	startStr := start.Format("2006-01-02")
	endStr := end.Format("2006-01-02")

	journals, err := svcCtx.Repos.Summary.FindByPeriodRangeAndSource(ctx, userID, constvar.SummaryPeriodTypeDay, startStr, endStr, constvar.SummarySourceUser)
	if err != nil || len(journals) == 0 {
		return "（无用户记录）"
	}

	var sb strings.Builder
	for _, j := range journals {
		content := j.SummaryContent
		if len([]rune(content)) > 300 {
			content = string([]rune(content)[:300]) + "..."
		}
		sb.WriteString(fmt.Sprintf("【%s】%s\n", j.PeriodStart, content))
	}
	return sb.String()
}

type monthlyLifeSignals struct {
	JournalCount int
	ExpenseTotal float64
	TopTags      []string
}

// buildLifetimeContext 人生总结额外上下文：标签统计 + 长期趋势
func buildLifetimeContext(ctx context.Context, svcCtx *svc.ServiceContext, userID uint64, start, end time.Time) string {
	startStr := start.Format("2006-01-02")
	endStr := end.Format("2006-01-02")

	tags, _ := svcCtx.Repos.Summary.ListTagsByDateRange(ctx, userID, startStr, endStr)
	tagPeriods, _ := svcCtx.Repos.Summary.ListTagPeriods(ctx, userID, startStr, endStr)
	journals, _ := svcCtx.Repos.Summary.FindByPeriodRangeAndSource(ctx, userID, constvar.SummaryPeriodTypeDay, startStr, endStr, constvar.SummarySourceUser)
	expenseMonths, _ := svcCtx.Repos.Expense.SumByMonth(ctx, userID, start, end)

	countMap, totalTags := accumulateTagCounts(tags)
	if len(countMap) == 0 && len(journals) == 0 && len(expenseMonths) == 0 {
		return ""
	}

	type kv struct {
		tag   string
		count int64
	}
	var sorted []kv
	for t, c := range countMap {
		sorted = append(sorted, kv{t, c})
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].count > sorted[j].count })

	var sb strings.Builder
	if len(sorted) > 0 && totalTags > 0 {
		sb.WriteString("\n【长期标签偏好】\n")
		for _, kv := range sorted {
			pct := float64(kv.count) / float64(totalTags) * 100
			sb.WriteString(fmt.Sprintf("  %s：%d 次 (%.1f%%)\n", kv.tag, kv.count, pct))
		}
	}

	monthlySignals := buildMonthlyLifeSignals(journals, tagPeriods, expenseMonths)
	if monthlyTrend := formatMonthlyLifeSignals(monthlySignals); monthlyTrend != "" {
		sb.WriteString("\n【长期行为趋势】\n")
		sb.WriteString(monthlyTrend)
	}

	return sb.String()
}

func accumulateTagCounts(tagStrings []string) (map[string]int64, int64) {
	countMap := make(map[string]int64)
	var total int64
	for _, tagStr := range tagStrings {
		for _, tag := range strings.Split(tagStr, ",") {
			tag = strings.TrimSpace(tag)
			if tag == "" {
				continue
			}
			countMap[tag]++
			total++
		}
	}
	return countMap, total
}

func buildMonthlyLifeSignals(journals []*summary.Summary, tagPeriods []summary.TagPeriod, expenseMonths []expenseRepo.MonthTotal) map[string]*monthlyLifeSignals {
	monthly := make(map[string]*monthlyLifeSignals)
	for _, journal := range journals {
		month := monthKeyFromDate(journal.PeriodStart)
		if month == "" {
			continue
		}
		if monthly[month] == nil {
			monthly[month] = &monthlyLifeSignals{}
		}
		monthly[month].JournalCount++
	}

	for _, monthTotal := range expenseMonths {
		if monthly[monthTotal.Month] == nil {
			monthly[monthTotal.Month] = &monthlyLifeSignals{}
		}
		monthly[monthTotal.Month].ExpenseTotal = monthTotal.Total
	}

	for _, period := range tagPeriods {
		if monthly[period.Month] == nil {
			monthly[period.Month] = &monthlyLifeSignals{}
		}
		countMap, _ := accumulateTagCounts([]string{period.Tags})
		for tag, count := range countMap {
			for i := int64(0); i < count; i++ {
				monthly[period.Month].TopTags = append(monthly[period.Month].TopTags, tag)
			}
		}
	}

	for month, signal := range monthly {
		if len(signal.TopTags) == 0 {
			continue
		}
		tagCounts, _ := accumulateTagCounts([]string{strings.Join(signal.TopTags, ",")})
		type kv struct {
			tag   string
			count int64
		}
		var sorted []kv
		for tag, count := range tagCounts {
			sorted = append(sorted, kv{tag: tag, count: count})
		}
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].count > sorted[j].count })
		limit := minInt(3, len(sorted))
		topTags := make([]string, 0, limit)
		for i := 0; i < limit; i++ {
			topTags = append(topTags, fmt.Sprintf("%s(%d)", sorted[i].tag, sorted[i].count))
		}
		monthly[month].TopTags = topTags
	}

	return monthly
}

func formatMonthlyLifeSignals(monthly map[string]*monthlyLifeSignals) string {
	if len(monthly) == 0 {
		return ""
	}

	var months []string
	for month := range monthly {
		months = append(months, month)
	}
	sort.Strings(months)
	if len(months) > 12 {
		months = months[len(months)-12:]
	}

	var sb strings.Builder
	for _, month := range months {
		signal := monthly[month]
		sb.WriteString(fmt.Sprintf("  %s：记录 %d 条，支出 %.2f 元", month, signal.JournalCount, signal.ExpenseTotal))
		if len(signal.TopTags) > 0 {
			sb.WriteString("，高频标签：")
			sb.WriteString(strings.Join(signal.TopTags, "、"))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func monthKeyFromDate(date string) string {
	if len(date) < 7 {
		return ""
	}
	return date[:7]
}

// buildContext 构建下级总结上下文
// 日报无需上下文；周报聚合日报；月报聚合周报；年报聚合月报
func buildContext(ctx context.Context, svcCtx *svc.ServiceContext, periodType uint8, userID uint64, start, end time.Time) string {
	childType := childPeriodType(periodType)
	if childType == 0 {
		return "（无历史上下文）"
	}

	summaries, err := svcCtx.Repos.Summary.FindByPeriodRange(ctx, userID, childType, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil || len(summaries) == 0 {
		return "（无历史上下文）"
	}

	var sb strings.Builder
	childLabel := periodTypeLabelCN(childType)
	for _, s := range summaries {
		// 截取前 200 字避免上下文过长
		content := s.SummaryContent
		if len([]rune(content)) > 200 {
			content = string([]rune(content)[:200]) + "..."
		}
		sb.WriteString(fmt.Sprintf("【%s %s】%s\n", childLabel, s.PeriodStart, content))
	}
	return sb.String()
}

// calcPeriod 计算当前周期（日报=昨天，周/月/年报=上个完整周期）
func calcPeriod(periodType uint8) (start, end time.Time) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	switch periodType {
	case constvar.SummaryPeriodTypeDay:
		start = today.AddDate(0, 0, -1) // 昨天
		end = today
	case constvar.SummaryPeriodTypeWeek:
		weekday := int(today.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		thisMonday := today.AddDate(0, 0, -weekday+1)
		start = thisMonday.AddDate(0, 0, -7) // 上周一
		end = thisMonday
	case constvar.SummaryPeriodTypeMonth:
		thisMonthFirst := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		start = thisMonthFirst.AddDate(0, -1, 0) // 上月 1 号
		end = thisMonthFirst
	case constvar.SummaryPeriodTypeYear:
		thisYearFirst := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		start = thisYearFirst.AddDate(-1, 0, 0) // 去年 1 月 1 号
		end = thisYearFirst
	case constvar.SummaryPeriodTypeLife:
		start = time.Date(2000, 1, 1, 0, 0, 0, 0, now.Location())
		end = now
	}
	return
}

// calcPeriodEnd 根据起始日期和周期类型计算结束日期
func calcPeriodEnd(periodType uint8, start time.Time) time.Time {
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	switch periodType {
	case constvar.SummaryPeriodTypeDay:
		return start.AddDate(0, 0, 1)
	case constvar.SummaryPeriodTypeWeek:
		return start.AddDate(0, 0, 7)
	case constvar.SummaryPeriodTypeMonth:
		return start.AddDate(0, 1, 0)
	case constvar.SummaryPeriodTypeYear:
		return start.AddDate(1, 0, 0)
	case constvar.SummaryPeriodTypeLife:
		return time.Now()
	default:
		return start.AddDate(0, 0, 1)
	}
}

// childPeriodType 返回下级周期类型（周报的下级是日报，月报的下级是周报，年报的下级是月报，人生总结的下级是年报）
func childPeriodType(t uint8) uint8 {
	switch t {
	case constvar.SummaryPeriodTypeWeek:
		return constvar.SummaryPeriodTypeDay
	case constvar.SummaryPeriodTypeMonth:
		return constvar.SummaryPeriodTypeWeek
	case constvar.SummaryPeriodTypeYear:
		return constvar.SummaryPeriodTypeMonth
	case constvar.SummaryPeriodTypeLife:
		return constvar.SummaryPeriodTypeYear
	default:
		return 0
	}
}

func formatCategoryBreakdown(ct []expenseRepo.CategoryTotal) string {
	if len(ct) == 0 {
		return "  无支出记录\n"
	}
	var sb strings.Builder
	for _, c := range ct {
		sb.WriteString(fmt.Sprintf("  %s：%.2f 元\n", c.CategoryName, c.Total))
	}
	return sb.String()
}

func periodTypeLabelCN(t uint8) string {
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

// buildLocationBreakdown 按地点汇总支出
func buildLocationBreakdown(ctx context.Context, svcCtx *svc.ServiceContext, userID uint64, start, end time.Time) string {
	logs, err := svcCtx.Repos.Expense.ListLogsByDateRange(ctx, userID, start, end)
	if err != nil || len(logs) == 0 {
		return "  无地点记录\n"
	}

	// 按地点汇总金额
	locMap := make(map[string]float64)
	for _, l := range logs {
		loc := l.Location
		if loc == "" {
			loc = "未知"
		}
		locMap[loc] += l.Amount
	}

	var sb strings.Builder
	for loc, total := range locMap {
		sb.WriteString(fmt.Sprintf("  %s：%.2f 元\n", loc, total))
	}
	return sb.String()
}
