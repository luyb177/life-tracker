<script setup lang="ts">
import { NTag, NSpace, NButton, NPopconfirm } from 'naive-ui'
import type { SummaryInfo } from '@/entities/summary/api/summary.api'
import { periodLabel, parseTags, formatPeriodRange } from '@/shared/utils/format'

defineProps<{ item: SummaryInfo; showActions?: boolean }>()
const emit = defineEmits<{ edit: [id: number]; delete: [id: number] }>()
</script>
<template>
  <div style="background: var(--n-color); border-radius: 8px; padding: 16px; margin-bottom: 12px; border: 1px solid var(--n-border-color)">
    <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 8px">
      <div>
        <NTag :type="item.source === 1 ? 'info' : 'success'" size="small" style="margin-right: 8px">{{ periodLabel(item.period_type) }}</NTag>
        <NTag :bordered="false" size="small">{{ item.source === 1 ? 'AI' : '用户' }}</NTag>
        <span style="font-size: 12px; color: var(--n-text-color-3); margin-left: 8px">{{ formatPeriodRange(item.period_start, item.period_end) }}</span>
      </div>
      <NSpace v-if="showActions" size="small">
        <NButton text size="tiny" @click="emit('edit', item.id)">编辑</NButton>
        <NPopconfirm @positive-click="emit('delete', item.id)"><template #trigger><NButton text size="tiny" type="error">删</NButton></template>确定删除？</NPopconfirm>
      </NSpace>
    </div>
    <div style="white-space: pre-wrap; line-height: 1.6; font-size: 14px">{{ item.summary_content }}</div>
    <NSpace v-if="parseTags(item.tags).length" style="margin-top: 8px">
      <NTag v-for="t in parseTags(item.tags)" :key="t" size="tiny" :bordered="false">{{ t }}</NTag>
    </NSpace>
    <div v-if="item.location && item.location !== '未知'" style="margin-top: 8px; font-size: 12px; color: var(--n-text-color-3)">📍 {{ item.location }}</div>
  </div>
</template>
