<script setup lang="ts">
import { NButton, NSpace } from 'naive-ui'
import EmptyState from '@/shared/ui/empty-state/EmptyState.vue'
import ConfirmAction from '@/shared/ui/confirm-action/ConfirmAction.vue'
import type { ExpenseLogInfo } from '@/entities/expense/api/expense.api'

defineProps<{ list: ExpenseLogInfo[]; total: number }>()
const emit = defineEmits<{ edit: [item: ExpenseLogInfo]; delete: [id: number] }>()
</script>
<template>
  <template v-if="list.length">
    <div v-for="e in list" :key="e.id" style="display: flex; justify-content: space-between; align-items: center; padding: 10px 0; border-bottom: 1px solid var(--n-border-color)">
      <div>
        <div style="font-weight: 500">{{ e.category?.name || '?' }} <span style="font-weight: 700; font-size: 16px; margin-left: 8px">¥{{ e.amount.toFixed(2) }}</span></div>
        <div v-if="e.note" style="font-size: 12px; color: var(--n-text-color-3)">{{ e.note }}</div>
      </div>
      <NSpace size="small">
        <NButton text size="tiny" @click="emit('edit', e)">编辑</NButton>
        <ConfirmAction @confirm="emit('delete', e.id)" />
      </NSpace>
    </div>
    <div style="text-align: right; padding-top: 8px; font-weight: 700">合计: ¥{{ total.toFixed(2) }}</div>
  </template>
  <EmptyState v-else title="当天暂无支出" />
</template>
