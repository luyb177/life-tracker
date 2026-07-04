<template>
  <div ref="chartRef" class="chart-canvas" />
</template>

<script setup lang="ts">
import { BarChart, LineChart } from 'echarts/charts'
import {
  GridComponent,
  TooltipComponent,
  type GridComponentOption,
  type TooltipComponentOption
} from 'echarts/components'
import * as echarts from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

echarts.use([BarChart, LineChart, GridComponent, TooltipComponent, CanvasRenderer])

export interface ExpenseTrendChartPoint {
  label: string
  total: number
}

const props = defineProps<{
  points: ExpenseTrendChartPoint[]
  type?: 'line' | 'bar'
}>()

const chartRef = ref<HTMLDivElement | null>(null)
let chart: echarts.ECharts | null = null

const chartType = computed(() => props.type || 'line')

function formatYuan(value: number) {
  return `¥${(value / 100).toFixed(2)}`
}

function renderChart() {
  if (!chartRef.value) return
  if (!chart) chart = echarts.init(chartRef.value)

  const labels = props.points.map((item) => item.label)
  const values = props.points.map((item) => item.total)
  const hasLongLabels = labels.some((label) => label.length > 10)

  const option: echarts.ComposeOption<GridComponentOption | TooltipComponentOption> = {
    grid: { top: 20, right: 12, bottom: 32, left: 52 },
    tooltip: {
      trigger: 'axis',
      valueFormatter: (value) => formatYuan(Number(value))
    },
    xAxis: {
      type: 'category',
      boundaryGap: chartType.value === 'bar',
      data: labels,
      axisTick: { show: false },
      axisLine: { lineStyle: { color: '#dfe5df' } },
      axisLabel: {
        color: '#66736f',
        rotate: hasLongLabels ? 30 : 0,
        hideOverlap: true
      }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        color: '#66736f',
        formatter: (value: number) => `¥${value / 100}`
      },
      splitLine: { lineStyle: { color: '#edf0ed' } }
    },
    series: [
      {
        type: chartType.value,
        data: values,
        smooth: chartType.value === 'line',
        symbolSize: 7,
        barMaxWidth: 26,
        itemStyle: { color: chartType.value === 'bar' ? '#d8664f' : '#168a7a' },
        lineStyle: { color: '#168a7a', width: 3 },
        areaStyle:
          chartType.value === 'line'
            ? { color: 'rgba(22, 138, 122, 0.12)' }
            : undefined
      }
    ]
  }

  chart.setOption(option, true)
}

function resizeChart() {
  chart?.resize()
}

watch(
  () => [props.points, props.type],
  async () => {
    await nextTick()
    renderChart()
  },
  { deep: true }
)

onMounted(async () => {
  await nextTick()
  renderChart()
  window.addEventListener('resize', resizeChart)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', resizeChart)
  chart?.dispose()
  chart = null
})
</script>
