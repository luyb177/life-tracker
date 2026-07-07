<template>
  <div class="page">
    <header class="page-header today-header">
      <div class="today-header-copy">
        <h1>{{ dateTitle }}</h1>
        <p class="page-description">随手记录生活与支出，再把一天复盘清楚。</p>
      </div>
      <div class="page-header-actions today-date-action">
        <n-date-picker
          class="today-date-picker"
          v-model:formatted-value="selectedDate"
          value-format="yyyy-MM-dd"
          type="date"
          :is-date-disabled="isFutureDateTimestamp"
        />
      </div>
    </header>

    <div class="metric-grid">
      <MetricCard label="本日支出" :value="formatYuan(totalExpense)" hint="正常支出合计" :icon="Wallet" tone="coral" />
      <MetricCard label="生活记录" :value="`${lifeLogs.length} 条`" :hint="`${dateTitle}已记录`" :icon="BookOpenText" tone="teal" />
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
      <TimelineList :items="timelineItems" @refund="refundTodayExpense" @select="openTimelineDetail" />
      <SummaryPreview
        :summary="dailySummary"
        :loading="generating"
        @generate="generateSummary"
        @open="openSummaryDetail"
      />
    </div>

    <LifeLogDetailModal
      v-model:show="showLifeModal"
      :log="selectedLifeLog"
      :loading="savingDetail"
      @save="saveLifeLog"
      @delete="deleteLifeLogByID"
    />
    <ExpenseDetailModal
      v-model:show="showExpenseModal"
      :expense="selectedExpense"
      :categories="expenseCategories"
      :loading="savingDetail"
      :loading-categories="loadingCategories"
      @save="saveExpense"
      @delete="deleteExpenseByID"
      @refund="refundSelectedExpense"
    />
    <SummaryDetailModal
      v-model:show="showSummaryModal"
      :summary="selectedSummary"
      :loading="savingDetail"
      @save="saveSummary"
      @delete="deleteSummaryByID"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useMessage } from 'naive-ui'
import { BookOpenText, Sparkles, Wallet } from '@lucide/vue'
import MetricCard from '@/components/common/MetricCard.vue'
import QuickLifeLogForm from '@/components/records/QuickLifeLogForm.vue'
import QuickExpenseForm from '@/components/records/QuickExpenseForm.vue'
import TimelineList, { type TimelineItem } from '@/components/records/TimelineList.vue'
import ExpenseDetailModal from '@/components/records/ExpenseDetailModal.vue'
import LifeLogDetailModal from '@/components/records/LifeLogDetailModal.vue'
import SummaryDetailModal from '@/components/summary/SummaryDetailModal.vue'
import SummaryPreview from '@/components/summary/SummaryPreview.vue'
import { deleteLifeLog, getLifeLogsByDate, updateLifeLog } from '@/api/lifeLog'
import {
  deleteExpense,
  getExpensesByDate,
  listExpenseCategories,
  refundExpense,
  updateExpense
} from '@/api/expense'
import { deleteSummary, generateAISummary, rangeSummaries, updateSummary } from '@/api/summary'
import type { ExpenseCategoryInfo, ExpenseLogInfo, LifeLogInfo, SummaryInfo } from '@/types/api'
import { addDays, formatDate, isFutureDateTimestamp, relativeDateLabel } from '@/utils/date'
import { formatYuan } from '@/utils/money'

const message = useMessage()
const selectedDate = ref(formatDate())
const lifeLogs = ref<LifeLogInfo[]>([])
const expenses = ref<ExpenseLogInfo[]>([])
const totalExpense = ref(0)
const summaries = ref<SummaryInfo[]>([])
const loading = ref(false)
const generating = ref(false)
const savingDetail = ref(false)
const showLifeModal = ref(false)
const showExpenseModal = ref(false)
const showSummaryModal = ref(false)
const selectedLifeLog = ref<LifeLogInfo | null>(null)
const selectedExpense = ref<ExpenseLogInfo | null>(null)
const selectedSummary = ref<SummaryInfo | null>(null)
const expenseCategories = ref<ExpenseCategoryInfo[]>([])
const loadingCategories = ref(false)
const dateTitle = computed(() => relativeDateLabel(selectedDate.value))

const dailySummary = computed(
  () => summaries.value.find((item) => item.period_type === 1 && item.source === 1) || null
)
const timelineItems = computed<TimelineItem[]>(() => {
  const lifeItems = lifeLogs.value.map((item) => ({
    id: `life-${item.id}`,
    type: 'life' as const,
    time: formatTimelineTime(item.occurred_at),
    sortAt: item.occurred_at,
    createdAt: item.created_at,
    sequence: item.id,
    title: '生活记录',
    description: item.content,
    location: item.location,
    tags: item.tags?.map((tag) => tag.name)
  }))
  const expenseItems = expenses.value.map((item) => ({
    id: `expense-${item.id}`,
    type: 'expense' as const,
    time: formatTimelineTime(item.occurred_at),
    sortAt: item.occurred_at,
    createdAt: item.created_at,
    sequence: item.id,
    title: item.category?.name || '支出',
    description: item.note || '未填写备注',
    location: item.location,
    amount: formatYuan(item.amount),
    canRefund: item.status === 0,
    refunded: item.status === 1
  }))
  return [...lifeItems, ...expenseItems].sort(compareTimelineItems)
})

function parseDateTime(value: string) {
  const timestamp = new Date(value.replace(' ', 'T')).getTime()
  return Number.isNaN(timestamp) ? 0 : timestamp
}

async function refundTodayExpense(id: number) {
  try {
    await refundExpense(id)
    message.success('退款成功')
    await loadToday()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '退款失败')
  }
}

function openTimelineDetail(item: TimelineItem) {
  if (item.type === 'life') {
    const log = lifeLogs.value.find((entry) => entry.id === item.sequence)
    if (!log) return
    selectedLifeLog.value = log
    showLifeModal.value = true
    return
  }

  const expense = expenses.value.find((entry) => entry.id === item.sequence)
  if (!expense) return
  selectedExpense.value = expense
  void loadExpenseCategories()
  showExpenseModal.value = true
}

function openSummaryDetail(summary: SummaryInfo) {
  selectedSummary.value = summary
  showSummaryModal.value = true
}

async function loadExpenseCategories() {
  if (expenseCategories.value.length > 0) return
  loadingCategories.value = true
  try {
    const resp = await listExpenseCategories()
    expenseCategories.value = resp.categories
  } finally {
    loadingCategories.value = false
  }
}

async function saveLifeLog(payload: { id: number; content: string; occurred_at: string }) {
  savingDetail.value = true
  try {
    await updateLifeLog(payload)
    message.success('生活记录已更新')
    showLifeModal.value = false
    await loadToday()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '更新失败')
  } finally {
    savingDetail.value = false
  }
}

async function deleteLifeLogByID(id: number) {
  if (!window.confirm('确认删除这条生活记录吗？')) return
  savingDetail.value = true
  try {
    await deleteLifeLog(id)
    message.success('生活记录已删除')
    showLifeModal.value = false
    await loadToday()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '删除失败')
  } finally {
    savingDetail.value = false
  }
}

async function saveExpense(payload: { id: number; category_id: number; amount: number; note: string; occurred_at: string }) {
  savingDetail.value = true
  try {
    await updateExpense(payload)
    message.success('支出已更新')
    showExpenseModal.value = false
    await loadToday()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '更新失败')
  } finally {
    savingDetail.value = false
  }
}

async function refundSelectedExpense(id: number) {
  await refundTodayExpense(id)
  showExpenseModal.value = false
}

async function deleteExpenseByID(id: number) {
  if (!window.confirm('确认删除这条支出记录吗？')) return
  savingDetail.value = true
  try {
    await deleteExpense(id)
    message.success('支出已删除')
    showExpenseModal.value = false
    await loadToday()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '删除失败')
  } finally {
    savingDetail.value = false
  }
}

async function saveSummary(payload: { id: number; title: string; summary_content: string; suggestion_content: string }) {
  savingDetail.value = true
  try {
    await updateSummary(payload)
    message.success('总结已更新')
    showSummaryModal.value = false
    await loadToday()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '更新失败')
  } finally {
    savingDetail.value = false
  }
}

async function deleteSummaryByID(id: number) {
  if (!window.confirm('确认删除这条总结吗？')) return
  savingDetail.value = true
  try {
    await deleteSummary(id)
    message.success('总结已删除')
    showSummaryModal.value = false
    await loadToday()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '删除失败')
  } finally {
    savingDetail.value = false
  }
}

function formatTimelineTime(value: string) {
  return value.slice(11, 19)
}

function compareTimelineItems(a: TimelineItem, b: TimelineItem) {
  const occurredDiff = parseDateTime(b.sortAt) - parseDateTime(a.sortAt)
  if (occurredDiff !== 0) return occurredDiff

  const createdDiff = parseDateTime(b.createdAt) - parseDateTime(a.createdAt)
  if (createdDiff !== 0) return createdDiff

  return b.sequence - a.sequence
}

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
