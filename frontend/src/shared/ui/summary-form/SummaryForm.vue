<script setup lang="ts">
import { NInput, NButton } from 'naive-ui'
import { useSummaryForm } from '@/shared/composables/useSummaryForm'
import PeriodInput from '@/shared/ui/period-input/PeriodInput.vue'
import TagInput from '@/shared/ui/tag-input/TagInput.vue'

const props = defineProps<{
  editId?: number | null
  initialContent?: string
  initialTags?: string
  initialTitle?: string
  initialPeriodType?: number
  initialPeriodStart?: string
  initialPeriodEnd?: string
  submitting?: boolean
  showPeriodSelect?: boolean
  lockPeriodFields?: boolean
}>()

const emit = defineEmits<{
  submit: [data: { period_type: number; period_start: string; period_end: string; summary_content: string; tags: string; title?: string }]
  cancel: []
}>()

const {
  periodType,
  ts,
  periodStart,
  content,
  title,
  tags,
  isLife,
  lifeEndTs,
  lifeError,
  periodLocked,
  displayPeriodEnd,
  displayPeriodLabel,
  ptOptions,
  buildPayload,
} = useSummaryForm({
  initialContent: props.initialContent,
  initialTags: props.initialTags,
  initialTitle: props.initialTitle,
  initialPeriodType: props.initialPeriodType,
  initialPeriodStart: props.initialPeriodStart,
  initialPeriodEnd: props.initialPeriodEnd,
  lockPeriodFields: props.lockPeriodFields,
})

function handleSubmit() {
  const payload = buildPayload()
  if (payload) emit('submit', payload)
}
</script>
<template>
  <div>
    <PeriodInput
      v-if="showPeriodSelect"
      :period-type="periodType"
      :period-start="periodStart"
      :period-end="displayPeriodEnd"
      :display-period-label="displayPeriodLabel"
      :period-locked="periodLocked"
      :is-life="isLife"
      :life-error="lifeError"
      :period-options="ptOptions"
      :start-value="ts"
      :end-value="lifeEndTs"
      @update:period-type="periodType = $event"
      @update:start-value="ts = $event"
      @update:end-value="lifeEndTs = $event"
    />
    <NInput v-model:value="title" placeholder="标题（可选）" style="margin-bottom: 12px" />
    <NInput v-model:value="content" type="textarea" placeholder="写点什么..." :autosize="{ minRows: 3, maxRows: 10 }" style="margin-bottom: 12px" />
    <TagInput v-model:model-value="tags" />
    <div style="display: flex; gap: 8px; margin-top: 12px; justify-content: flex-end">
      <NButton @click="emit('cancel')">取消</NButton>
      <NButton type="primary" :loading="submitting" :disabled="!periodLocked && lifeError" @click="handleSubmit">{{ editId ? '保存' : '提交' }}</NButton>
    </div>
  </div>
</template>
