<script setup lang="ts">
import { NModal } from 'naive-ui'
import type { SummaryInfo } from '@/entities/summary/api/summary.api'
import SummaryForm from '@/shared/ui/summary-form/SummaryForm.vue'

const props = defineProps<{
  show: boolean
  title: string
  submitting?: boolean
  item?: SummaryInfo | null
  showPeriodSelect?: boolean
  lockPeriodFields?: boolean
  initialPeriodType?: number
  initialPeriodStart?: string
  initialPeriodEnd?: string
}>()

const emit = defineEmits<{
  'update:show': [value: boolean]
  submit: [data: { period_type: number; period_start: string; period_end: string; summary_content: string; tags: string; title?: string }]
}>()

function close() {
  emit('update:show', false)
}
</script>

<template>
  <NModal :show="show" :title="title" @update:show="emit('update:show', $event)">
    <div style="padding: 16px">
      <SummaryForm
        :edit-id="item?.id"
        :initial-content="item?.summary_content"
        :initial-tags="item?.tags"
        :initial-title="item?.title"
        :initial-period-type="item?.period_type ?? initialPeriodType"
        :initial-period-start="item?.period_start ?? initialPeriodStart"
        :initial-period-end="item?.period_end ?? initialPeriodEnd"
        :show-period-select="showPeriodSelect"
        :lock-period-fields="lockPeriodFields"
        :submitting="submitting"
        @submit="emit('submit', $event)"
        @cancel="close"
      />
    </div>
  </NModal>
</template>
