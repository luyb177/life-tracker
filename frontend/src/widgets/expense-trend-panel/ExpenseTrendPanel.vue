<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, GridComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import ChartCard from '@/shared/ui/chart-card/ChartCard.vue'
import { getStatsTrend } from '@/entities/expense/api/expense.api'
import { expenseKeys } from '@/entities/expense/api/expense.keys'

use([LineChart, TitleComponent, TooltipComponent, GridComponent, CanvasRenderer])

const props = defineProps<{ start: string; end: string; title?: string }>()

const { data, isLoading } = useQuery({
  queryKey: computed(() => expenseKeys.statsTrend(props.start, props.end)),
  queryFn: () => getStatsTrend(props.start, props.end),
})

const points = computed(() => (data.value as any)?.points || [])

const option = computed(() => ({
  tooltip: { trigger: 'axis' as const },
  grid: { left: 40, right: 16, top: 16, bottom: 24 },
  xAxis: { type: 'category' as const, data: points.value.map((p: any) => p.date.slice(5)) },
  yAxis: { type: 'value' as const },
  series: [{ type: 'line' as const, data: points.value.map((p: any) => p.total), smooth: true, areaStyle: { opacity: 0.1 } }],
}))
</script>
<template>
  <ChartCard :title="title || '支出趋势'" :loading="isLoading" :is-empty="!points.length">
    <VChart :option="option" style="height: 220px" />
  </ChartCard>
</template>
