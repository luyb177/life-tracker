<script setup lang="ts">
import { NDatePicker, NSpace, NSelect } from 'naive-ui'
import { PERIOD_LABELS } from '@/entities/summary/api/summary.api'

const props = defineProps<{
  date?: number | null
  dateRange?: [number, number] | null
  periodType?: number | null
  showPeriod?: boolean
  showDate?: boolean
  showRange?: boolean
}>()
const emit = defineEmits<{
  'update:date': [v: number]
  'update:dateRange': [v: [number, number]]
  'update:periodType': [v: number]
}>()

const ptOptions = Object.entries(PERIOD_LABELS).map(([k, v]) => ({ label: v, value: Number(k) }))
</script>
<template>
  <NSpace align="center" style="margin-bottom: 16px" :wrap="true">
    <NDatePicker v-if="showDate" :value="props.date" type="date" size="small" style="width: 160px" @update:value="emit('update:date', $event as number)" />
    <NDatePicker v-if="showRange" :value="props.dateRange" type="daterange" size="small" style="width: 240px" @update:value="emit('update:dateRange', $event as [number, number])" />
    <NSelect v-if="showPeriod" :value="props.periodType" :options="ptOptions" size="small" style="width: 100px" @update:value="emit('update:periodType', $event)" />
    <slot />
  </NSpace>
</template>
