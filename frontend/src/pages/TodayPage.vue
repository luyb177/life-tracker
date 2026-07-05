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

    <n-modal v-model:show="showLifeModal" preset="dialog" title="生活记录详情" class="detail-modal">
      <n-form class="detail-form" :show-label="true" label-placement="top">
        <n-form-item label="发生时间">
          <n-date-picker
            v-model:formatted-value="lifeEdit.occurred_at"
            type="datetime"
            value-format="yyyy-MM-dd HH:mm:ss"
            class="full-control"
          />
        </n-form-item>
        <n-form-item label="内容">
          <n-input v-model:value="lifeEdit.content" type="textarea" :autosize="{ minRows: 5, maxRows: 9 }" />
        </n-form-item>
        <p v-if="selectedLifeLog" class="detail-meta">
          创建 {{ selectedLifeLog.created_at }} · 更新 {{ selectedLifeLog.last_updated_at || selectedLifeLog.updated_at }}
        </p>
      </n-form>
      <template #action>
        <n-button tertiary type="error" :loading="savingDetail" @click="deleteSelectedLifeLog">删除</n-button>
        <n-button @click="showLifeModal = false">取消</n-button>
        <n-button type="primary" :loading="savingDetail" @click="saveSelectedLifeLog">保存</n-button>
      </template>
    </n-modal>

    <n-modal v-model:show="showExpenseModal" preset="dialog" title="支出详情" class="detail-modal">
      <n-form class="detail-form" :show-label="true" label-placement="top">
        <n-form-item label="发生时间">
          <n-date-picker
            v-model:formatted-value="expenseEdit.occurred_at"
            type="datetime"
            value-format="yyyy-MM-dd HH:mm:ss"
            class="full-control"
          />
        </n-form-item>
        <div class="detail-two-columns">
          <n-form-item label="金额">
            <n-input-number v-model:value="expenseEdit.amount" :min="0" :precision="2" class="full-control">
              <template #prefix>¥</template>
            </n-input-number>
          </n-form-item>
          <n-form-item label="分类">
            <n-select
              v-model:value="expenseEdit.category_id"
              :options="categoryOptions"
              :loading="loadingCategories"
              class="full-control"
            />
          </n-form-item>
        </div>
        <n-form-item label="备注">
          <n-input v-model:value="expenseEdit.note" />
        </n-form-item>
        <p v-if="selectedExpense" class="detail-meta">
          创建 {{ selectedExpense.created_at }} · 更新 {{ selectedExpense.last_updated_at || selectedExpense.updated_at }}
        </p>
      </n-form>
      <template #action>
        <n-button tertiary type="error" :loading="savingDetail" @click="deleteSelectedExpense">删除</n-button>
        <n-button v-if="selectedExpense?.status === 0" tertiary type="warning" :loading="savingDetail" @click="refundSelectedExpense">
          退款
        </n-button>
        <n-button @click="showExpenseModal = false">取消</n-button>
        <n-button type="primary" :loading="savingDetail" @click="saveSelectedExpense">保存</n-button>
      </template>
    </n-modal>

    <n-modal v-model:show="showSummaryModal" preset="dialog" title="总结详情" class="detail-modal">
      <n-form class="detail-form" :show-label="true" label-placement="top">
        <n-form-item label="标题">
          <n-input v-model:value="summaryEdit.title" />
        </n-form-item>
        <n-form-item label="总结内容">
          <n-input v-model:value="summaryEdit.summary_content" type="textarea" :autosize="{ minRows: 6, maxRows: 12 }" />
        </n-form-item>
        <n-form-item label="建议">
          <n-input v-model:value="summaryEdit.suggestion_content" type="textarea" :autosize="{ minRows: 3, maxRows: 8 }" />
        </n-form-item>
        <p v-if="selectedSummary" class="detail-meta">
          来源 {{ selectedSummary.source === 1 ? 'AI' : '用户' }} · 更新 {{ selectedSummary.last_updated_at || selectedSummary.updated_at }}
        </p>
      </n-form>
      <template #action>
        <n-button tertiary type="error" :loading="savingDetail" @click="deleteSelectedSummary">删除</n-button>
        <n-button @click="showSummaryModal = false">取消</n-button>
        <n-button type="primary" :loading="savingDetail" @click="saveSelectedSummary">保存</n-button>
      </template>
    </n-modal>
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
import { fenToYuan, formatYuan, yuanToFen } from '@/utils/money'

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
const lifeEdit = ref({ content: '', occurred_at: '' })
const expenseEdit = ref({ amount: 0, category_id: null as number | null, note: '', occurred_at: '' })
const summaryEdit = ref({ title: '', summary_content: '', suggestion_content: '' })
const dateTitle = computed(() => relativeDateLabel(selectedDate.value))

const dailySummary = computed(
  () => summaries.value.find((item) => item.period_type === 1 && item.source === 1) || null
)
const categoryOptions = computed(() =>
  expenseCategories.value.map((item) => ({ label: item.name, value: item.id }))
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
    lifeEdit.value = {
      content: log.content,
      occurred_at: log.occurred_at
    }
    showLifeModal.value = true
    return
  }

  const expense = expenses.value.find((entry) => entry.id === item.sequence)
  if (!expense) return
  selectedExpense.value = expense
  expenseEdit.value = {
    amount: fenToYuan(expense.amount),
    category_id: expense.category.id,
    note: expense.note,
    occurred_at: expense.occurred_at
  }
  void loadExpenseCategories()
  showExpenseModal.value = true
}

function openSummaryDetail(summary: SummaryInfo) {
  selectedSummary.value = summary
  summaryEdit.value = {
    title: summary.title || '',
    summary_content: summary.summary_content,
    suggestion_content: summary.suggestion_content || ''
  }
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

async function saveSelectedLifeLog() {
  if (!selectedLifeLog.value || !lifeEdit.value.content.trim()) return
  savingDetail.value = true
  try {
    await updateLifeLog({
      id: selectedLifeLog.value.id,
      content: lifeEdit.value.content.trim(),
      occurred_at: lifeEdit.value.occurred_at
    })
    message.success('生活记录已更新')
    showLifeModal.value = false
    await loadToday()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '更新失败')
  } finally {
    savingDetail.value = false
  }
}

async function deleteSelectedLifeLog() {
  if (!selectedLifeLog.value || !window.confirm('确认删除这条生活记录吗？')) return
  savingDetail.value = true
  try {
    await deleteLifeLog(selectedLifeLog.value.id)
    message.success('生活记录已删除')
    showLifeModal.value = false
    await loadToday()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '删除失败')
  } finally {
    savingDetail.value = false
  }
}

async function saveSelectedExpense() {
  if (!selectedExpense.value || !expenseEdit.value.category_id || expenseEdit.value.amount <= 0) return
  savingDetail.value = true
  try {
    await updateExpense({
      id: selectedExpense.value.id,
      category_id: expenseEdit.value.category_id,
      amount: yuanToFen(expenseEdit.value.amount),
      note: expenseEdit.value.note.trim(),
      occurred_at: expenseEdit.value.occurred_at
    })
    message.success('支出已更新')
    showExpenseModal.value = false
    await loadToday()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '更新失败')
  } finally {
    savingDetail.value = false
  }
}

async function refundSelectedExpense() {
  if (!selectedExpense.value) return
  await refundTodayExpense(selectedExpense.value.id)
  showExpenseModal.value = false
}

async function deleteSelectedExpense() {
  if (!selectedExpense.value || !window.confirm('确认删除这条支出记录吗？')) return
  savingDetail.value = true
  try {
    await deleteExpense(selectedExpense.value.id)
    message.success('支出已删除')
    showExpenseModal.value = false
    await loadToday()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '删除失败')
  } finally {
    savingDetail.value = false
  }
}

async function saveSelectedSummary() {
  if (!selectedSummary.value || !summaryEdit.value.summary_content.trim()) return
  savingDetail.value = true
  try {
    await updateSummary({
      id: selectedSummary.value.id,
      title: summaryEdit.value.title.trim(),
      summary_content: summaryEdit.value.summary_content.trim(),
      suggestion_content: summaryEdit.value.suggestion_content.trim()
    })
    message.success('总结已更新')
    showSummaryModal.value = false
    await loadToday()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '更新失败')
  } finally {
    savingDetail.value = false
  }
}

async function deleteSelectedSummary() {
  if (!selectedSummary.value || !window.confirm('确认删除这条总结吗？')) return
  savingDetail.value = true
  try {
    await deleteSummary(selectedSummary.value.id)
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
