<template>
  <div class="page">
    <PageHeader title="总结" description="查看日、周、月、年的 AI 和手动复盘。">
      <n-button type="primary" secondary :loading="generating" @click="generateCurrent">
        生成当前周期
      </n-button>
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
              <strong>{{ summary.title || periodLabel(summary.period_type) }}</strong>
              <p>{{ summary.period_start }} 至 {{ summary.period_end }}</p>
            </div>
            <span class="source-badge">{{ summary.source === 1 ? 'AI' : '用户' }}</span>
          </div>
          <p>{{ summary.summary_content }}</p>
          <p v-if="summary.suggestion_content" class="suggestion">{{ summary.suggestion_content }}</p>
        </article>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useMessage } from 'naive-ui'
import { Sparkles } from '@lucide/vue'
import PageHeader from '@/components/common/PageHeader.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import { generateAISummary, listSummaries } from '@/api/summary'
import type { SummaryInfo } from '@/types/api'
import { formatDate, monthStart } from '@/utils/date'

const message = useMessage()
const periodType = ref(1)
const summaries = ref<SummaryInfo[]>([])
const loading = ref(false)
const generating = ref(false)

function periodLabel(type: number) {
  return ['未知', '日报', '周报', '月报', '年报'][type] || '总结'
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

watch(periodType, load)
onMounted(load)
</script>

