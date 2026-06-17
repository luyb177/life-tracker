<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { NSpin } from 'naive-ui'
import StatCard from '@/shared/ui/stat-card/StatCard.vue'
import EmptyState from '@/shared/ui/empty-state/EmptyState.vue'
import { getExpenseByDate, ExpenseLogInfo } from '@/entities/expense/api/expense.api'
import { expenseKeys } from '@/entities/expense/api/expense.keys'

const props = defineProps<{ date: string }>()

const { data, isLoading } = useQuery({
  queryKey: computed(() => expenseKeys.byDate(props.date)),
  queryFn: () => getExpenseByDate(props.date),
  enabled: computed(() => !!props.date),
})

const list = computed(() => ((data.value as any)?.list || []) as ExpenseLogInfo[])
const total = computed(() => (data.value as any)?.total || 0)
</script>
<template>
  <NSpin :show="isLoading">
    <div style="display: flex; gap: 12px; margin-bottom: 12px">
      <StatCard title="今日支出" :value="total.toFixed(2)" unit="元" />
      <StatCard title="笔数" :value="list.length" unit="笔" />
    </div>
    <template v-if="list.length">
      <div v-for="e in list" :key="e.id" style="display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid var(--n-border-color)">
        <span>{{ e.category?.name }}<span v-if="e.note" style="color: var(--n-text-color-3); margin-left: 4px">- {{ e.note }}</span></span>
        <span style="font-weight: 600">¥{{ e.amount.toFixed(2) }}</span>
      </div>
      <div style="text-align: right; padding-top: 8px; font-weight: 700">合计: ¥{{ total.toFixed(2) }}</div>
    </template>
    <EmptyState v-else title="暂无支出" />
  </NSpin>
</template>
