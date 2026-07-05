<template>
  <n-modal :show="show" preset="dialog" title="总结详情" class="detail-modal" @update:show="emit('update:show', $event)">
    <n-form class="detail-form" :show-label="true" label-placement="top">
      <n-form-item label="标题">
        <n-input v-model:value="form.title" />
      </n-form-item>
      <n-form-item label="总结内容">
        <n-input v-model:value="form.summary_content" type="textarea" :autosize="{ minRows: 6, maxRows: 12 }" />
      </n-form-item>
      <n-form-item label="建议">
        <n-input v-model:value="form.suggestion_content" type="textarea" :autosize="{ minRows: 3, maxRows: 8 }" />
      </n-form-item>
      <p v-if="summary" class="detail-meta">
        {{ periodLabel(summary.period_type) }} · {{ summary.source === 1 ? 'AI' : '用户' }} · 更新
        {{ summary.last_updated_at || summary.updated_at }}
      </p>
    </n-form>
    <template #action>
      <n-button tertiary type="error" :loading="loading" :disabled="!summary" @click="remove">删除</n-button>
      <n-button @click="emit('update:show', false)">取消</n-button>
      <n-button type="primary" :loading="loading" :disabled="!canSave" @click="save">保存</n-button>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { SummaryInfo } from '@/types/api'
import { periodLabel } from '@/utils/summary'

const props = defineProps<{ show: boolean; summary?: SummaryInfo | null; loading?: boolean }>()
const emit = defineEmits<{
  'update:show': [value: boolean]
  save: [payload: { id: number; title: string; summary_content: string; suggestion_content: string }]
  delete: [id: number]
}>()

const form = ref({ title: '', summary_content: '', suggestion_content: '' })
const canSave = computed(() => Boolean(props.summary && form.value.summary_content.trim()))

watch(
  () => [props.summary, props.show] as const,
  ([summary, show]) => {
    if (!summary || !show) return
    form.value = {
      title: summary.title || '',
      summary_content: summary.summary_content,
      suggestion_content: summary.suggestion_content || ''
    }
  },
  { immediate: true }
)

function save() {
  if (!props.summary || !canSave.value) return
  emit('save', {
    id: props.summary.id,
    title: form.value.title.trim(),
    summary_content: form.value.summary_content.trim(),
    suggestion_content: form.value.suggestion_content.trim()
  })
}

function remove() {
  if (!props.summary) return
  emit('delete', props.summary.id)
}
</script>
