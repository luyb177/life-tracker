<script setup lang="ts">
import { NAlert, NDatePicker, NSelect } from 'naive-ui'
import { formatPeriodRange } from '@/shared/utils/format'

defineProps<{
  periodType: number
  periodStart: string
  periodEnd: string
  displayPeriodLabel: string
  periodLocked?: boolean
  isLife?: boolean
  lifeError?: boolean
  periodOptions: Array<{ label: string; value: number }>
  startValue: number
  endValue?: number
}>()

const emit = defineEmits<{
  'update:periodType': [value: number]
  'update:startValue': [value: number]
  'update:endValue': [value: number]
}>()
</script>
<template>
  <div style="margin-bottom: 12px">
    <template v-if="periodLocked">
      <div style="display: flex; gap: 8px; align-items: center; flex-wrap: wrap; margin-bottom: 8px">
        <span style="font-size: 13px; font-weight: 600">{{ displayPeriodLabel }}</span>
        <span style="font-size: 12px; color: var(--n-text-color-3)">{{ formatPeriodRange(periodStart, periodEnd) }}</span>
      </div>
      <NAlert type="info" size="small">编辑时仅支持修改标题、内容和标签，周期范围保持不变。</NAlert>
    </template>
    <template v-else>
      <NSelect
        :value="periodType"
        :options="periodOptions"
        size="small"
        style="width: 120px; margin-bottom: 8px"
        @update:value="emit('update:periodType', $event)"
      />
      <div style="display: flex; gap: 8px; align-items: center; flex-wrap: wrap">
        <NDatePicker
          :value="startValue"
          type="date"
          size="small"
          style="width: 160px"
          @update:value="emit('update:startValue', $event as number)"
        />
        <span style="font-size: 12px; color: var(--n-text-color-3)">→ {{ periodEnd }}</span>
      </div>
      <div v-if="isLife" style="margin-top: 8px; display: flex; gap: 8px; align-items: center">
        <span style="font-size: 12px; color: var(--n-text-color-3)">结束:</span>
        <NDatePicker
          :value="endValue"
          type="date"
          size="small"
          style="width: 160px"
          @update:value="emit('update:endValue', $event as number)"
        />
      </div>
    </template>
    <NAlert v-if="lifeError" type="error" size="small" style="margin-top: 8px">结束日期必须晚于起始日期</NAlert>
  </div>
</template>
