<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { PieChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

use([PieChart, TitleComponent, TooltipComponent, LegendComponent, CanvasRenderer])

const props = defineProps<{ data: any[]; title?: string }>()
const option = computed(() => ({
  tooltip: { trigger: 'item' as const },
  legend: { bottom: 0 },
  series: [{ type: 'pie' as const, radius: ['40%', '70%'], center: ['50%', '45%'], data: props.data }],
}))
</script>
<template>
  <div style="margin-bottom: 24px">
    <h4 style="margin: 0 0 8px">{{ title || '分类占比' }}</h4>
    <VChart v-if="data.length" :option="option" style="height: 260px" />
  </div>
</template>
