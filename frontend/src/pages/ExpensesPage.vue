<template>
  <div class="page">
    <PageHeader title="支出" description="查看每日支出、分类结构和近期趋势。">
      <n-date-picker
        v-model:formatted-value="date"
        value-format="yyyy-MM-dd"
        type="date"
        :is-date-disabled="isFutureDateTimestamp"
      />
    </PageHeader>

    <div class="metric-grid">
      <MetricCard label="当日合计" :value="formatYuan(total)" hint="已退款不计入" :icon="Wallet" tone="coral" />
      <MetricCard label="支出笔数" :value="`${expenses.length} 笔`" hint="选中日期" :icon="ReceiptText" tone="teal" />
      <MetricCard label="分类数量" :value="`${categories.length} 个`" hint="系统与自定义" :icon="Tags" tone="amber" />
    </div>

    <section class="panel expense-workspace">
      <n-tabs v-model:value="activeTab" type="segment" animated>
        <n-tab-pane name="records" tab="明细">
          <div class="section-title-row">
            <div>
              <p class="eyebrow">Records</p>
              <h2>{{ date }} 的支出</h2>
            </div>
            <n-button secondary :loading="loading" @click="load">刷新</n-button>
          </div>

          <EmptyState
            v-if="expenses.length === 0"
            title="这天还没有支出"
            description="回到今天页面可以快速记一笔花销。"
            :icon="ReceiptText"
          />
          <div v-else class="record-list expense-record-list">
            <article
              v-for="item in expenses"
              :key="item.id"
              class="record-row split interactive-record"
              :class="{ refunded: item.status === 1 }"
              @click="openExpenseDetail(item)"
            >
              <div>
                <div class="record-meta-row">
                  <time>{{ item.occurred_at }}</time>
                  <span v-if="item.location" class="record-location-inline">地点：{{ item.location }}</span>
                </div>
                <p>{{ item.category.name }} · {{ item.note || '无备注' }}</p>
                <span v-if="item.status === 1" class="status-badge">已退款</span>
              </div>
              <div class="expense-record-actions">
                <strong>{{ formatYuan(item.amount) }}</strong>
                <n-button
                  v-if="item.status === 0"
                  size="tiny"
                  tertiary
                  type="error"
                  @click.stop="refundRecord(item.id)"
                >
                  退款
                </n-button>
              </div>
            </article>
          </div>
        </n-tab-pane>

        <n-tab-pane name="trend" tab="趋势">
          <div class="section-title-row">
            <div>
              <p class="eyebrow">Trend</p>
              <h2>{{ trendTitle }}</h2>
            </div>
            <n-button secondary :loading="trendLoading" @click="loadTrend">刷新</n-button>
          </div>

          <n-tabs v-model:value="trendMode" type="segment" size="small" class="nested-tabs">
            <n-tab name="today">今日</n-tab>
            <n-tab name="month">本月</n-tab>
            <n-tab name="year">今年</n-tab>
          </n-tabs>

          <EmptyState
            v-if="trendPoints.length === 0"
            title="暂无趋势数据"
            description="有支出记录后，这里会展示花销变化。"
            :icon="ReceiptText"
          />
          <div v-else>
            <ExpenseTrendChart :points="trendPoints" :type="trendChartType" />
            <p class="chart-caption">{{ trendCaption }}</p>
          </div>
        </n-tab-pane>

        <n-tab-pane name="categories" tab="分类">
          <div class="section-title-row">
            <div>
              <p class="eyebrow">Categories</p>
              <h2>支出分类</h2>
            </div>
          </div>
          <div class="category-list">
            <span v-for="category in categories" :key="category.id" class="chip readonly">
              {{ category.name }}
            </span>
            <button type="button" class="chip add-chip" @click="openCategoryModal">
              + 添加分类
            </button>
          </div>
        </n-tab-pane>
      </n-tabs>
    </section>

    <n-modal v-model:show="showCategoryModal" preset="dialog" title="添加支出分类">
      <n-input
        v-model:value="categoryName"
        placeholder="例如：咖啡、打车、吹头发"
        maxlength="50"
        show-count
        @keyup.enter="createCategory"
      />
      <template #action>
        <n-button @click="showCategoryModal = false">取消</n-button>
        <n-button type="primary" :loading="creatingCategory" @click="createCategory">保存</n-button>
      </template>
    </n-modal>

    <ExpenseDetailModal
      v-model:show="showExpenseModal"
      :expense="selectedExpense"
      :categories="categories"
      :loading="savingExpense"
      @save="saveSelectedExpense"
      @delete="deleteSelectedExpense"
      @refund="refundSelectedExpense"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useMessage } from 'naive-ui'
import { ReceiptText, Tags, Wallet } from '@lucide/vue'
import PageHeader from '@/components/common/PageHeader.vue'
import MetricCard from '@/components/common/MetricCard.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import ExpenseTrendChart, { type ExpenseTrendChartPoint } from '@/components/charts/ExpenseTrendChart.vue'
import ExpenseDetailModal from '@/components/records/ExpenseDetailModal.vue'
import {
  createExpenseCategory,
  deleteExpense,
  getExpensesByDate,
  getMonthlyTrendStats,
  getTrendStats,
  listExpenseCategories,
  refundExpense,
  updateExpense
} from '@/api/expense'
import type { ExpenseCategoryInfo, ExpenseLogInfo } from '@/types/api'
import {
  formatDate,
  isFutureDateTimestamp,
  monthEndExclusive,
  monthStart,
  nextYearStart,
  yearStart
} from '@/utils/date'
import { formatYuan } from '@/utils/money'

const message = useMessage()
const date = ref(formatDate())
const activeTab = ref<'records' | 'trend' | 'categories'>('records')
const trendMode = ref<'today' | 'month' | 'year'>('today')
const expenses = ref<ExpenseLogInfo[]>([])
const categories = ref<ExpenseCategoryInfo[]>([])
const total = ref(0)
const loading = ref(false)
const trendLoading = ref(false)
const trendPoints = ref<ExpenseTrendChartPoint[]>([])
const showCategoryModal = ref(false)
const showExpenseModal = ref(false)
const creatingCategory = ref(false)
const savingExpense = ref(false)
const categoryName = ref('')
const selectedExpense = ref<ExpenseLogInfo | null>(null)

const trendTitle = computed(() => {
  if (trendMode.value === 'today') return `${date.value} 今日花销趋势`
  if (trendMode.value === 'month') return '本月花销趋势'
  return '今年花销趋势'
})

const trendCaption = computed(() => {
  if (trendMode.value === 'today') return '按小时聚合当天正常支出。'
  if (trendMode.value === 'month') return '按天聚合本月正常支出。'
  return '按月聚合今年正常支出。'
})

const trendChartType = computed<'line' | 'bar'>(() => (trendMode.value === 'year' ? 'bar' : 'line'))

async function load() {
  loading.value = true
  try {
    const [expenseResp, categoryResp] = await Promise.all([
      getExpensesByDate(date.value),
      listExpenseCategories()
    ])
    expenses.value = expenseResp.list
    total.value = expenseResp.total
    categories.value = categoryResp.categories
    if (trendMode.value === 'today') {
      trendPoints.value = buildTodayTrend(expenseResp.list)
    }
  } catch (error) {
    message.error(error instanceof Error ? error.message : '加载支出失败')
  } finally {
    loading.value = false
  }
}

function buildTodayTrend(list: ExpenseLogInfo[]) {
  const buckets = new Map<string, number>()
  list
    .filter((item) => item.status === 0)
    .forEach((item) => {
      const hour = item.occurred_at.slice(11, 13)
      const label = `${hour}:00`
      buckets.set(label, (buckets.get(label) || 0) + item.amount)
    })

  return Array.from(buckets.entries())
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([label, value]) => ({ label, total: value }))
}

async function loadTrend() {
  trendLoading.value = true
  try {
    if (trendMode.value === 'today') {
      trendPoints.value = buildTodayTrend(expenses.value)
      return
    }

    if (trendMode.value === 'month') {
      const resp = await getTrendStats({
        start: monthStart(new Date(`${date.value}T00:00:00`)),
        end: monthEndExclusive(new Date(`${date.value}T00:00:00`))
      })
      trendPoints.value = resp.points.map((item) => ({
        label: formatChartDateLabel(item.date),
        total: item.total
      }))
      return
    }

    const resp = await getMonthlyTrendStats({
      start: yearStart(new Date(`${date.value}T00:00:00`)),
      end: nextYearStart(new Date(`${date.value}T00:00:00`))
    })
    trendPoints.value = resp.points.map((item) => ({
      label: item.month.slice(0, 7),
      total: item.total
    }))
  } catch (error) {
    message.error(error instanceof Error ? error.message : '加载趋势失败')
  } finally {
    trendLoading.value = false
  }
}

function formatChartDateLabel(value: string) {
  return value.slice(0, 10)
}

function openCategoryModal() {
  categoryName.value = ''
  showCategoryModal.value = true
}

function openExpenseDetail(item: ExpenseLogInfo) {
  selectedExpense.value = item
  showExpenseModal.value = true
}

async function createCategory() {
  const name = categoryName.value.trim()
  if (!name) return
  creatingCategory.value = true
  try {
    await createExpenseCategory(name)
    const categoryResp = await listExpenseCategories()
    categories.value = categoryResp.categories
    showCategoryModal.value = false
    message.success('分类已添加')
  } catch (error) {
    message.error(error instanceof Error ? error.message : '添加分类失败')
  } finally {
    creatingCategory.value = false
  }
}

async function refundRecord(id: number) {
  try {
    await refundExpense(id)
    message.success('退款成功')
    await load()
    if (activeTab.value === 'trend') await loadTrend()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '退款失败')
  }
}

async function saveSelectedExpense(payload: { id: number; category_id: number; amount: number; note: string; occurred_at: string }) {
  savingExpense.value = true
  try {
    await updateExpense(payload)
    message.success('支出已更新')
    showExpenseModal.value = false
    await load()
    if (activeTab.value === 'trend') await loadTrend()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '更新失败')
  } finally {
    savingExpense.value = false
  }
}

async function refundSelectedExpense(id: number) {
  await refundRecord(id)
  showExpenseModal.value = false
}

async function deleteSelectedExpense(id: number) {
  if (!window.confirm('确认删除这条支出记录吗？')) return
  savingExpense.value = true
  try {
    await deleteExpense(id)
    message.success('支出已删除')
    showExpenseModal.value = false
    await load()
    if (activeTab.value === 'trend') await loadTrend()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '删除失败')
  } finally {
    savingExpense.value = false
  }
}

watch(date, load)
watch([trendMode, activeTab], () => {
  if (activeTab.value === 'trend') loadTrend()
})
onMounted(load)
</script>
