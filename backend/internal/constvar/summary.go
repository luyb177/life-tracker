package constvar

// 总结周期类型
const (
	SummaryPeriodTypeDay   uint8 = 1 // 日报
	SummaryPeriodTypeWeek  uint8 = 2 // 周报
	SummaryPeriodTypeMonth uint8 = 3 // 月报
	SummaryPeriodTypeYear  uint8 = 4 // 年报
)

// 总结来源
const (
	SummarySourceAI   uint8 = 1 // AI 生成
	SummarySourceUser uint8 = 2 // 用户手动
)

// 总结分页游标
const (
	SummaryPageTokenPrefix = "summary"
)

// 支出分类
const (
	ExpenseCategoryTypeSystem uint8 = 1 // 系统默认
	ExpenseCategoryTypeUser   uint8 = 2 // 用户自定义
)

// 默认分类名称
var DefaultExpenseCategories = []string{"早饭", "午饭", "晚饭", "杂项"}

// 支出分页游标
const (
	ExpensePageTokenPrefix = "expense"
)

// 生活记录分页游标
const (
	LifeLogPageTokenPrefix = "lifelog"
)
