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
        <article
          v-for="summary in summaries"
          :key="summary.id"
          class="summary-card interactive-summary"
          @click="openSummaryDetail(summary)"
        >
          <div class="summary-card-header">
            <div>
              <strong class="summary-period">{{ summaryPeriodTitle(summary) }}</strong>
              <p v-if="summary.title" class="summary-subtitle">{{ summary.title }}</p>
              <p v-else class="summary-subtitle">{{ summaryPeriodRange(summary) }}</p>
              <p v-if="summaryLocationText(summary)" class="record-location">地点：{{ summaryLocationText(summary) }}</p>
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

    <SummaryDetailModal
      v-model:show="showDetailModal"
      :summary="selectedSummary"
      :loading="savingDetail"
      @save="saveSelectedSummary"
      @delete="deleteSelectedSummary"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useMessage } from 'naive-ui'
import { Sparkles } from '@lucide/vue'
import PageHeader from '@/components/common/PageHeader.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import SummaryDetailModal from '@/components/summary/SummaryDetailModal.vue'
import { createSummary, deleteSummary, generateAISummary, listSummaries, updateSummary } from '@/api/summary'
import type { SummaryInfo } from '@/types/api'
import { formatDate, isFutureDateTimestamp } from '@/utils/date'
import { renderMarkdown } from '@/utils/markdown'
import {
  currentPeriodStart,
  normalizePeriodStart,
  periodEnd,
  summaryLocationText,
  summaryPeriodOptions,
  summaryPeriodRange,
  summaryPeriodTitle
} from '@/utils/summary'

const message = useMessage()
const periodType = ref(1)
const summaries = ref<SummaryInfo[]>([])
const loading = ref(false)
const generating = ref(false)
const showManualModal = ref(false)
const showDetailModal = ref(false)
const creatingManual = ref(false)
const savingDetail = ref(false)
const selectedSummary = ref<SummaryInfo | null>(null)
const manualPeriodType = ref(1)
const manualPeriodStart = ref(formatDate())
const manualTitle = ref('')
const manualContent = ref('')
const manualSuggestion = ref('')
const periodOptions = summaryPeriodOptions

function renderSummaryMarkdown(content = '') {
  return renderMarkdown(content)
}

function openManualSummary() {
  manualPeriodType.value = periodType.value
  manualPeriodStart.value = currentPeriodStart(periodType.value)
  manualTitle.value = ''
  manualContent.value = ''
  manualSuggestion.value = ''
  showManualModal.value = true
}

function openSummaryDetail(summary: SummaryInfo) {
  selectedSummary.value = summary
  showDetailModal.value = true
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
    const summary = await generateAISummary(periodType.value, currentPeriodStart(periodType.value))
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
    const normalizedStart = normalizePeriodStart(manualPeriodType.value, manualPeriodStart.value)
    await createSummary({
      period_type: manualPeriodType.value,
      period_start: normalizedStart,
      period_end: periodEnd(manualPeriodType.value, normalizedStart),
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

async function saveSelectedSummary(payload: { id: number; title: string; summary_content: string; suggestion_content: string }) {
  savingDetail.value = true
  try {
    await updateSummary(payload)
    message.success('总结已更新')
    showDetailModal.value = false
    await load()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '更新失败')
  } finally {
    savingDetail.value = false
  }
}

async function deleteSelectedSummary(id: number) {
  if (!window.confirm('确认删除这条总结吗？')) return
  savingDetail.value = true
  try {
    await deleteSummary(id)
    message.success('总结已删除')
    showDetailModal.value = false
    await load()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '删除失败')
  } finally {
    savingDetail.value = false
  }
}

watch(periodType, load)
onMounted(load)
</script>
