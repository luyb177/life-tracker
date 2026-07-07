<template>
  <section class="panel summary-preview" :class="{ interactive: Boolean(summary) }" @click="summary && $emit('open', summary)">
    <div class="section-title-row">
      <div>
        <p class="eyebrow">AI</p>
        <h2>今日总结</h2>
      </div>
      <n-button size="small" secondary :loading="loading" @click.stop="$emit('generate')">
        生成
      </n-button>
    </div>

    <EmptyState
      v-if="!summary"
      title="暂无 AI 总结"
      description="记录一些内容后，可以手动生成今日复盘。"
      :icon="Sparkles"
    />
    <div v-else class="summary-content">
      <strong>{{ summary.title || '今日复盘' }}</strong>
      <p v-if="locationText" class="record-location">地点：{{ locationText }}</p>
      <div class="markdown-body" v-html="summaryHtml" />
      <div v-if="suggestionHtml" class="markdown-body suggestion" v-html="suggestionHtml" />
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Sparkles } from '@lucide/vue'
import EmptyState from '@/components/common/EmptyState.vue'
import type { SummaryInfo } from '@/types/api'
import { renderMarkdown } from '@/utils/markdown'
import { summaryLocationText } from '@/utils/summary'

const props = defineProps<{ summary?: SummaryInfo | null; loading?: boolean }>()
defineEmits<{ generate: []; open: [summary: SummaryInfo] }>()

const summaryHtml = computed(() => renderMarkdown(props.summary?.summary_content || ''))
const suggestionHtml = computed(() => renderMarkdown(props.summary?.suggestion_content || ''))
const locationText = computed(() => summaryLocationText(props.summary))
</script>
