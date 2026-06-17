<script setup lang="ts">
import { ref, computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import FilterBar from '@/shared/ui/filter-bar/FilterBar.vue'
import PageHeader from '@/shared/ui/page-header/PageHeader.vue'
import PageSection from '@/shared/ui/page-section/PageSection.vue'
import ExpenseTrendPanel from '@/widgets/expense-trend-panel/ExpenseTrendPanel.vue'
import LifeSummaryPanel from '@/widgets/life-summary-panel/LifeSummaryPanel.vue'
import CategoryBreakdownWidget from '@/widgets/expense-trend-panel/CategoryBreakdownWidget.vue'
import TagTrendPanel from '@/widgets/tag-trend-panel/TagTrendPanel.vue'
import AnalyticsOverviewPanel from '@/widgets/analytics-panels/AnalyticsOverviewPanel.vue'
import MonthlyExpenseTrendPanel from '@/widgets/analytics-panels/MonthlyExpenseTrendPanel.vue'
import TagStatsPanel from '@/widgets/analytics-panels/TagStatsPanel.vue'
import { getStatsCategory } from '@/entities/expense/api/expense.api'
import { expenseKeys } from '@/entities/expense/api/expense.keys'
import { formatDate } from '@/shared/utils/format'

const range = ref<[number, number]>([Date.now() - 30 * 86400000, Date.now()])
const start = computed(() => formatDate(new Date(range.value[0])))
const end = computed(() => formatDate(new Date(range.value[1])))

const { data: category } = useQuery({
  queryKey: computed(() => expenseKeys.statsCategory(start.value, end.value)),
  queryFn: () => getStatsCategory(start.value, end.value),
  enabled: computed(() => !!start.value),
})

const pieData = computed(() => ((category.value as any)?.categories || []).map((c: any) => ({ name: c.category_name, value: c.total })))
</script>
<template>
  <PageHeader title="分析" description="支出趋势与标签统计" />
  <FilterBar v-model:date-range="range" show-range />
  <PageSection title="区间概览" description="先看这段时间的支出总量和标签活跃度。">
    <AnalyticsOverviewPanel :start="start" :end="end" />
  </PageSection>
  <PageSection title="趋势变化" description="从日趋势和月趋势判断开销变化是否稳定。">
    <ExpenseTrendPanel :start="start" :end="end" title="每日支出趋势" />
    <MonthlyExpenseTrendPanel :start="start" :end="end" />
  </PageSection>
  <PageSection title="分类与标签" description="从支出分类和标签频次里总结最近的生活主题。">
    <CategoryBreakdownWidget :data="pieData" />
    <TagStatsPanel :start="start" :end="end" />
    <TagTrendPanel :start="start" :end="end" />
  </PageSection>
  <PageSection title="长期复盘" description="把周期数据与长期人生总结放到一起看。">
    <LifeSummaryPanel />
  </PageSection>
</template>
