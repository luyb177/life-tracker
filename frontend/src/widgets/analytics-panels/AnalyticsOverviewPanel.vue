<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import StatCard from '@/shared/ui/stat-card/StatCard.vue'
import { getStatsRange } from '@/entities/expense/api/expense.api'
import { getTagStats } from '@/entities/summary/api/summary.api'
import { expenseKeys } from '@/entities/expense/api/expense.keys'
import { summaryKeys } from '@/entities/summary/api/summary.keys'

const props = defineProps<{ start: string; end: string }>()

const { data: rangeTotal, isLoading: rangeLoading } = useQuery({
  queryKey: computed(() => expenseKeys.statsRange(props.start, props.end)),
  queryFn: () => getStatsRange(props.start, props.end),
  enabled: computed(() => !!props.start),
})

const { data: tags, isLoading: tagsLoading } = useQuery({
  queryKey: computed(() => summaryKeys.tagStats({ start: props.start, end: props.end })),
  queryFn: () => getTagStats({ start: props.start, end: props.end }),
  enabled: computed(() => !!props.start),
})

const total = computed(() => (rangeTotal.value as any)?.total ?? 0)
const tagCount = computed(() => (tags.value as any)?.tags?.length || 0)
</script>
<template>
  <div style="display: flex; gap: 12px; margin-bottom: 16px; flex-wrap: wrap">
   <StatCard title="区间总支出" :value="total.toFixed(2)" unit="元" :loading="rangeLoading" />
   <StatCard title="标签数" :value="tagCount" unit="个" :loading="tagsLoading" />
  </div>
</template>
