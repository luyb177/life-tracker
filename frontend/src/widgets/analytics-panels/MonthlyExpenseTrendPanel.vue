<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import ChartCard from '@/shared/ui/chart-card/ChartCard.vue'
import { getMonthlyTrend } from '@/entities/expense/api/expense.api'
import { expenseKeys } from '@/entities/expense/api/expense.keys'

const props = defineProps<{ start: string; end: string }>()

const { data, isLoading } = useQuery({
  queryKey: computed(() => expenseKeys.monthlyTrend(props.start, props.end)),
  queryFn: () => getMonthlyTrend(props.start, props.end),
  enabled: computed(() => !!props.start),
})

const points = computed(() => (data.value as any)?.points || [])
</script>
<template>
  <ChartCard title="每月支出趋势" :loading="isLoading" :is-empty="!points.length">
    <div v-if="points.length" style="display: flex; gap: 8px; flex-wrap: wrap">
      <div v-for="m in points" :key="m.month" style="background: var(--n-color-target); border-radius: 6px; padding: 8px 12px; min-width: 80px; text-align: center">
        <div style="font-size: 11px; color: var(--n-text-color-3)">{{ m.month }}</div>
        <div style="font-weight: 700; font-size: 16px">¥{{ m.total.toFixed(0) }}</div>
      </div>
    </div>
  </ChartCard>
</template>
