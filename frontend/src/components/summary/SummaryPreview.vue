<template>
  <section class="panel summary-preview">
    <div class="section-title-row">
      <div>
        <p class="eyebrow">AI</p>
        <h2>今日总结</h2>
      </div>
      <n-button size="small" secondary :loading="loading" @click="$emit('generate')">
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
      <p>{{ summary.summary_content }}</p>
      <p v-if="summary.suggestion_content" class="suggestion">{{ summary.suggestion_content }}</p>
    </div>
  </section>
</template>

<script setup lang="ts">
import { Sparkles } from '@lucide/vue'
import EmptyState from '@/components/common/EmptyState.vue'
import type { SummaryInfo } from '@/types/api'

defineProps<{ summary?: SummaryInfo | null; loading?: boolean }>()
defineEmits<{ generate: [] }>()
</script>

