<template>
  <div class="page">
    <PageHeader title="总结" description="查看日、周、月、年的 AI 和手动复盘。">
      <n-button secondary @click="openManualSummary">写手动总结</n-button>
      <n-button type="primary" secondary :loading="generating" @click="generateCurrent">生成当前周期</n-button>
    </PageHeader>

    <section class="panel">
      <div class="section-title-row">
        <n-tabs v-model:value="periodType" type="segment" size="small">
          <n-tab :name="1">日</n-tab>
          <n-tab :name="2">周</n-tab>
          <n-tab :name="3">月</n-tab>
          <n-tab :name="4">年</n-tab>
        </n-tabs>
        <n-button secondary :loading="loading" @click="load">刷新</n-button>
      </div>

      <EmptyState
        v-if="summaries.length === 0"
        title="暂无总结"
        description="可以先生成当前周期的 AI 总结。"
        :icon="Sparkles"
      />

      <div v-else class="summary-list">
        <article v-for="summary in summaries" :key="summary.id" class="summary-card">
          <div class="summary-card-header">
            <div>
              <strong class="summary-period">{{ summaryPeriodTitle(summary) }}</strong>
              <p v-if="summary.title" class="summary-subtitle">{{ summary.title }}</p>
              <p v-else class="summary-subtitle">{{ summaryPeriodRange(summary) }}</p>
            </div>
            <span class="source-badge">{{ summary.source === 1 ? 'AI' : '用户' }}</span>
          </div>
          <div class="markdown-body" v-html="renderSummaryMarkdown(summary.summary_content)" />
          <div
            v-if="summary.suggestion_content"
            class="markdown-body suggestion"
            v-html="renderSummaryMarkdown(summary.suggestion_content)"
          />
        </article>
      </div>
    </section>

    <n-modal v-model:show="showManualModal" preset="dialog" title="写手动总结" class="summary-modal">
      <n-form :show-label="true" label-placement="top">
        <n-form-item label="周期">
          <n-select v-model:value="manualPeriodType" :options="periodOptions" />
        </n-form-item>
        <n-form-item label="开始日期">
          <n-date-picker
            v-model:formatted-value="manualPeriodStart"
            value-format="yyyy-MM-dd"
            type="date"
            :is-date-disabled="isFutureDateTimestamp"
          />
        </n-form-item>
        <n-form-item label="标题">
          <n-input v-model:value="manualTitle" placeholder="可选" />
        </n-form-item>
        <n-form-item label="总结内容">
          <n-input
            v-model:value="manualContent"
            type="textarea"
            placeholder="支持 Markdown，例如 ## 今日复盘"
            :autosize="{ minRows: 6, maxRows: 12 }"
          />
        </n-form-item>
        <n-form-item label="建议">
          <n-input
            v-model:value="manualSuggestion"
            type="textarea"
            placeholder="可选，同样支持 Markdown"
            :autosize="{ minRows: 3, maxRows: 8 }"
          />
        </n-form-item>
      </n-form>
      <template #action>
        <n-button @click="showManualModal = false">取消</n-button>
        <n-button type="primary" :loading="creatingManual" :disabled="!manualContent.trim()" @click="createManualSummary">
          保存
        </n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useMessage } from 'naive-ui'
import { Sparkles } from '@lucide/vue'
import PageHeader from '@/components/common/PageHeader.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import { createSummary, generateAISummary, listSummaries } from '@/api/summary'
import type { SummaryInfo } from '@/types/api'
import { addDays, formatDate, isFutureDateTimestamp, monthEndExclusive, monthStart, nextYearStart } from '@/utils/date'
import { renderMarkdown } from '@/utils/markdown'

const message = useMessage()
const periodType = ref(1)
const summaries = ref<SummaryInfo[]>([])
const loading = ref(false)
const generating = ref(false)
const showManualModal = ref(false)
const creatingManual = ref(false)
const manualPeriodType = ref(1)
const manualPeriodStart = ref(formatDate())
const manualTitle = ref('')
const manualContent = ref('')
const manualSuggestion = ref('')

const periodOptions = [
  { label: '日', value: 1 },
  { label: '周', value: 2 },
  { label: '月', value: 3 },
  { label: '年', value: 4 }
]

function periodLabel(type: number) {
  return ['未知', '日报', '周报', '月报', '年报'][type] || '总结'
}

function renderSummaryMarkdown(content = '') {
  return renderMarkdown(content)
}

function summaryPeriodTitle(summary: SummaryInfo) {
  if (summary.period_type === 1) return summary.period_start
  if (summary.period_type === 3) return summary.period_start.slice(0, 7)
  if (summary.period_type === 4) return summary.period_start.slice(0, 4)
  return summaryPeriodRange(summary)
}

function summaryPeriodRange(summary: SummaryInfo) {
  return `${summary.period_start} 至 ${summary.period_end}`
}

function currentPeriodStart() {
  const today = new Date()
  if (periodType.value === 3) return monthStart(today)
  if (periodType.value === 4) return `${today.getFullYear()}-01-01`
  if (periodType.value === 2) {
    const day = today.getDay() || 7
    today.setDate(today.getDate() - day + 1)
  }
  return formatDate(today)
}

function normalizeManualStart(periodType: number, dateText: string) {
  const date = new Date(`${dateText}T00:00:00`)
  if (periodType === 2) {
    const day = date.getDay() || 7
    date.setDate(date.getDate() - day + 1)
  }
  if (periodType === 3) {
    date.setDate(1)
  }
  if (periodType === 4) {
    date.setMonth(0, 1)
  }
  return formatDate(date)
}

function manualPeriodEnd(periodType: number, start: string) {
  const startDate = new Date(`${start}T00:00:00`)
  if (periodType === 1) return addDays(start, 1)
  if (periodType === 2) return addDays(start, 7)
  if (periodType === 3) return monthEndExclusive(startDate)
  return nextYearStart(startDate)
}

function openManualSummary() {
  manualPeriodType.value = periodType.value
  manualPeriodStart.value = currentPeriodStart()
  manualTitle.value = ''
  manualContent.value = ''
  manualSuggestion.value = ''
  showManualModal.value = true
}

async function load() {
  loading.value = true
  try {
    summaries.value = (await listSummaries({ period_type: periodType.value, page_size: 20 })).list
  } catch (error) {
    message.error(error instanceof Error ? error.message : '加载总结失败')
  } finally {
    loading.value = false
  }
}

async function generateCurrent() {
  generating.value = true
  try {
    const summary = await generateAISummary(periodType.value, currentPeriodStart())
    summaries.value = [summary, ...summaries.value]
    message.success('总结已生成')
  } catch (error) {
    message.error(error instanceof Error ? error.message : '生成总结失败')
  } finally {
    generating.value = false
  }
}

async function createManualSummary() {
  const content = manualContent.value.trim()
  if (!content) return
  creatingManual.value = true
  try {
    const normalizedStart = normalizeManualStart(manualPeriodType.value, manualPeriodStart.value)
    await createSummary({
      period_type: manualPeriodType.value,
      period_start: normalizedStart,
      period_end: manualPeriodEnd(manualPeriodType.value, normalizedStart),
      title: manualTitle.value.trim() || undefined,
      summary_content: content,
      suggestion_content: manualSuggestion.value.trim() || undefined
    })
    showManualModal.value = false
    periodType.value = manualPeriodType.value
    await load()
    message.success('手动总结已保存')
  } catch (error) {
    message.error(error instanceof Error ? error.message : '保存手动总结失败')
  } finally {
    creatingManual.value = false
  }
}

watch(periodType, load)
onMounted(load)
</script>
