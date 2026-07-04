<template>
  <div class="page">
    <PageHeader title="今天" description="随手记录生活与支出，再把一天复盘清楚。">
      <n-date-picker v-model:formatted-value="selectedDate" value-format="yyyy-MM-dd" type="date" />
    </PageHeader>

    <div class="metric-grid">
      <MetricCard label="本日支出" :value="formatYuan(totalExpense)" hint="正常支出合计" :icon="Wallet" tone="coral" />
      <MetricCard label="生活记录" :value="`${lifeLogs.length} 条`" hint="今天已记录" :icon="BookOpenText" tone="teal" />
      <MetricCard label="AI 总结" :value="dailySummary ? '已生成' : '待生成'" hint="日报复盘" :icon="Sparkles" tone="amber" />
    </div>

    <div class="today-grid">
      <section class="panel form-panel">
        <QuickLifeLogForm :date="selectedDate" @created="loadToday" />
      </section>
      <section class="panel form-panel">
        <QuickExpenseForm :date="selectedDate" @created="loadToday" />
      </section>
    </div>

    <div class="today-grid wide-left">
      <TimelineList :items="timelineItems" />
      <SummaryPreview :summary="dailySummary" :loading="generating" @generate="generateSummary" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useMessage } from 'naive-ui'
import { BookOpenText, Sparkles, Wallet } from '@lucide/vue'
import PageHeader from '@/components/common/PageHeader.vue'
import MetricCard from '@/components/common/MetricCard.vue'
import QuickLifeLogForm from '@/components/records/QuickLifeLogForm.vue'
import QuickExpenseForm from '@/components/records/QuickExpenseForm.vue'
import TimelineList, { type TimelineItem } from '@/components/records/TimelineList.vue'
import SummaryPreview from '@/components/summary/SummaryPreview.vue'
import { getLifeLogsByDate } from '@/api/lifeLog'
import { getExpensesByDate } from '@/api/expense'
import { generateAISummary, rangeSummaries } from '@/api/summary'
import type { ExpenseLogInfo, LifeLogInfo, SummaryInfo } from '@/types/api'
import { addDays, formatDate } from '@/utils/date'
import { formatYuan } from '@/utils/money'

const message = useMessage()
const selectedDate = ref(formatDate())
const lifeLogs = ref<LifeLogInfo[]>([])
const expenses = ref<ExpenseLogInfo[]>([])
const totalExpense = ref(0)
const summaries = ref<SummaryInfo[]>([])
const loading = ref(false)
const generating = ref(false)

const dailySummary = computed(
  () => summaries.value.find((item) => item.period_type === 1 && item.source === 1) || null
)

const timelineItems = computed<TimelineItem[]>(() => {
  const lifeItems = lifeLogs.value.map((item) => ({
    id: `life-${item.id}`,
    type: 'life' as const,
    time: item.occurred_at.slice(11, 16),
    title: '生活记录',
    description: item.content,
    tags: item.tags?.map((tag) => tag.name)
  }))
  const expenseItems = expenses.value.map((item) => ({
    id: `expense-${item.id}`,
    type: 'expense' as const,
    time: item.occurred_at.slice(11, 16),
    title: item.category?.name || '支出',
    description: item.note || '未填写备注',
    amount: formatYuan(item.amount)
  }))
  return [...lifeItems, ...expenseItems].sort((a, b) => b.time.localeCompare(a.time))
})

async function loadToday() {
  loading.value = true
  try {
    const [lifeResp, expenseResp, summaryResp] = await Promise.all([
      getLifeLogsByDate(selectedDate.value),
      getExpensesByDate(selectedDate.value),
      rangeSummaries({ period_type: 1, start: selectedDate.value, end: addDays(selectedDate.value, 1) })
    ])
    lifeLogs.value = lifeResp.list
    expenses.value = expenseResp.list
    totalExpense.value = expenseResp.total
    summaries.value = summaryResp.list
  } catch (error) {
    message.error(error instanceof Error ? error.message : '加载今日数据失败')
  } finally {
    loading.value = false
  }
}

async function generateSummary() {
  generating.value = true
  try {
    const summary = await generateAISummary(1, selectedDate.value)
    summaries.value = [summary, ...summaries.value.filter((item) => item.id !== summary.id)]
    message.success('AI 总结已生成')
  } catch (error) {
    message.error(error instanceof Error ? error.message : '生成失败')
  } finally {
    generating.value = false
  }
}

watch(selectedDate, loadToday)
onMounted(loadToday)
</script>

