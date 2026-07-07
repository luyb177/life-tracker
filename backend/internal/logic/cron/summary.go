package summary

import (
	"context"
	"fmt"
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
// lastUpdatedBy 为空时表示系统任务更新，记为 0。
func Run(ctx context.Context, svcCtx *svc.ServiceContext, periodType uint8, userID uint64, periodStart time.Time, lastUpdatedBy ...uint64) error {
	var periodEnd time.Time
	if periodStart.IsZero() {
		periodStart, periodEnd = calcPeriod(periodType)
	} else {
		periodEnd = calcPeriodEnd(periodType, periodStart)
	}
	modifierID := uint64(0)
	if len(lastUpdatedBy) > 0 {
		modifierID = lastUpdatedBy[0]
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
	expenseDetailText := buildExpenseDetailContext(ctx, svcCtx, userID, periodStart, periodEnd)

	// 2. 构建上下文（下级总结）
	contextText := buildContext(ctx, svcCtx, periodType, userID, periodStart, periodEnd)

	// 2.5 获取用户日报（生活记录）
	journalText := buildJournalContext(ctx, svcCtx, userID, periodStart, periodEnd)

	// 3. 构建 Prompt
	label := periodTypeLabelCN(periodType)
	systemPrompt := "你是一个个人生活助手。请根据提供的数据生成简洁的「" + label + "」总结。用中文回复，语气亲切。"
	userPrompt := fmt.Sprintf(`请生成%s总结：

周期：%s ~ %s

【支出汇总】
总支出：%.2f 元
分类明细：
%s
【支出明细】
%s
【地点分布】
%s
【用户记录】
%s
【上下文】
%s

请包含：1. 消费概况与分类分析 2. 结合用户记录分析当日行为 3. 地点分布分析 4. 变化趋势（对比上下文） 5. 改进建议。
如果数据为空或极少，请如实说明"今日无事"或"本周期无记录"，不用强行编造。`,
		label,
		periodStart.Format("2006-01-02"), periodEnd.Format("2006-01-02"),
		float64(totalExpense)/100,
		formatCategoryBreakdown(categoryBreakdown),
		expenseDetailText,
		locationBreakdown,
		journalText,
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
	now := time.Now().In(constvar.TimeLocation)
	s := &summary.Summary{
		UserID:            userID,
		PeriodType:        periodType,
		PeriodStart:       periodStartStr,
		PeriodEnd:         periodEnd.Format("2006-01-02"),
		Source:            constvar.SummarySourceAI,
		SummaryContent:    aiContent,
		SuggestionContent: "",
		Location:          locationBreakdown,
		LastUpdatedBy:     modifierID,
		LastUpdatedAt:     now,
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
			"last_updated_by":    modifierID,
			"last_updated_at":    now,
		})
	}
	return svcCtx.Repos.Summary.Create(ctx, s)
}

// buildJournalContext 获取周期内的用户生活记录作为 AI 上下文
func buildJournalContext(ctx context.Context, svcCtx *svc.ServiceContext, userID uint64, start, end time.Time) string {
	logs, err := svcCtx.Repos.LifeLog.FindByDateRange(ctx, userID, start, end)
	if err != nil || len(logs) == 0 {
		return "（无用户记录）"
	}

	var sb strings.Builder
	for _, l := range logs {
		content := l.Content
		if len([]rune(content)) > 300 {
			content = string([]rune(content)[:300]) + "..."
		}
		location := strings.TrimSpace(l.Location)
		if location == "" {
			location = "未知"
		}
		sb.WriteString(fmt.Sprintf("【%s｜地点：%s】%s\n", l.OccurredAt.Format("2006-01-02 15:04"), location, content))
	}
	return sb.String()
}

// buildContext 构建下级总结上下文
// 日报无需上下文；周报聚合日报；月报聚合周报；年报聚合月报
func buildContext(ctx context.Context, svcCtx *svc.ServiceContext, periodType uint8, userID uint64, start, end time.Time) string {
	childType := childPeriodType(periodType)
	if childType == 0 {
		return "（无历史上下文）"
	}

	summaries, err := svcCtx.Repos.Summary.FindByPeriodRangeAndSource(ctx, userID, childType, start.Format("2006-01-02"), end.Format("2006-01-02"), constvar.SummarySourceAI)
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
	default:
		return start.AddDate(0, 0, 1)
	}
}

// childPeriodType 返回下级周期类型（周报的下级是日报，月报的下级是周报，年报的下级是月报）
func childPeriodType(t uint8) uint8 {
	switch t {
	case constvar.SummaryPeriodTypeWeek:
		return constvar.SummaryPeriodTypeDay
	case constvar.SummaryPeriodTypeMonth:
		return constvar.SummaryPeriodTypeWeek
	case constvar.SummaryPeriodTypeYear:
		return constvar.SummaryPeriodTypeMonth
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
		sb.WriteString(fmt.Sprintf("  %s：%.2f 元\n", c.CategoryName, float64(c.Total)/100))
	}
	return sb.String()
}

func buildExpenseDetailContext(ctx context.Context, svcCtx *svc.ServiceContext, userID uint64, start, end time.Time) string {
	logs, err := svcCtx.Repos.Expense.ListLogsByDateRange(ctx, userID, start, end)
	if err != nil || len(logs) == 0 {
		return "  无支出明细\n"
	}

	categories, err := svcCtx.Repos.Expense.FindCategoriesByUser(ctx, userID)
	categoryNames := make(map[uint64]string)
	if err == nil {
		for _, category := range categories {
			categoryNames[category.ID] = category.Name
		}
	}

	var sb strings.Builder
	count := 0
	for _, log := range logs {
		if log.Status == 1 {
			continue
		}
		count++
		categoryName := strings.TrimSpace(categoryNames[log.CategoryID])
		if categoryName == "" {
			categoryName = "未分类"
		}
		note := strings.TrimSpace(log.Note)
		if note == "" {
			note = "无备注"
		}
		if len([]rune(note)) > 120 {
			note = string([]rune(note)[:120]) + "..."
		}
		location := strings.TrimSpace(log.Location)
		if location == "" {
			location = "未知"
		}
		sb.WriteString(fmt.Sprintf(
			"  【%s】%s %.2f 元；备注：%s；地点：%s\n",
			log.OccurredAt.Format("2006-01-02 15:04"),
			categoryName,
			float64(log.Amount)/100,
			note,
			location,
		))
	}
	if count == 0 {
		return "  无支出明细\n"
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

	// 按地点汇总金额（排除已退款）
	locMap := make(map[string]int64)
	for _, l := range logs {
		if l.Status == 1 {
			continue
		}
		loc := l.Location
		if loc == "" {
			loc = "未知"
		}
		locMap[loc] += l.Amount
	}

	var sb strings.Builder
	for loc, total := range locMap {
		sb.WriteString(fmt.Sprintf("  %s：%.2f 元\n", loc, float64(total)/100))
	}
	return sb.String()
}
