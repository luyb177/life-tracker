<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'
import PageHeader from '@/shared/ui/page-header/PageHeader.vue'
import PageSection from '@/shared/ui/page-section/PageSection.vue'
import ExpenseTrendPanel from '@/widgets/expense-trend-panel/ExpenseTrendPanel.vue'
import DashboardOverviewPanel from '@/widgets/dashboard-overview/DashboardOverviewPanel.vue'
import DashboardDailyPanels from '@/widgets/dashboard-overview/DashboardDailyPanels.vue'
import { getExpenseByDate } from '@/entities/expense/api/expense.api'
import { getDaySummaries } from '@/entities/summary/api/summary.api'
import { expenseKeys } from '@/entities/expense/api/expense.keys'
import { summaryKeys } from '@/entities/summary/api/summary.keys'
import { todayStr, formatDate, isUserRecord } from '@/shared/utils/format'

const router = useRouter()
const today = todayStr()
const weekAgo = formatDate(new Date(Date.now() - 7 * 86400000))
const { data: dayData } = useQuery({ queryKey: summaryKeys.day(today), queryFn: () => getDaySummaries(today) })
const { data: expenseData } = useQuery({ queryKey: expenseKeys.byDate(today), queryFn: () => getExpenseByDate(today) })
const userRecords = computed(() => (((dayData.value as any)?.list || []) as any[]).filter(s => isUserRecord(s.source)))
const total = computed(() => (expenseData.value as any)?.total || 0)
const list = computed(() => (expenseData.value as any)?.list || [])
</script>
<template>
  <PageHeader title="仪表盘" :description="`${today}`" />
  <PageSection title="今日概览" description="快速查看今天的记录、支出和常用入口。">
    <DashboardOverviewPanel
      :record-count="userRecords.length"
      :total-expense="total"
      :expense-count="list.length"
      @records="router.push('/records')"
      @expenses="router.push('/expenses')"
      @summaries="router.push('/summaries')"
    />
  </PageSection>
  <PageSection title="趋势预览" description="先看最近一周的支出变化，再决定是否进入详细分析。">
    <ExpenseTrendPanel :start="weekAgo" :end="today" title="本周支出" />
  </PageSection>
  <PageSection title="今日内容" description="把今天的生活记录和支出明细放在同一个区域里查看。">
    <DashboardDailyPanels :date="today" />
  </PageSection>
</template>
