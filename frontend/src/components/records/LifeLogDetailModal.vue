<template>
  <n-modal :show="show" preset="dialog" title="生活记录详情" class="detail-modal" @update:show="emit('update:show', $event)">
    <n-form class="detail-form" :show-label="true" label-placement="top">
      <n-form-item label="发生时间">
        <n-date-picker
          v-model:formatted-value="form.occurred_at"
          type="datetime"
          value-format="yyyy-MM-dd HH:mm:ss"
          class="full-control"
        />
      </n-form-item>
      <n-form-item label="内容">
        <n-input v-model:value="form.content" type="textarea" :autosize="{ minRows: 5, maxRows: 9 }" />
      </n-form-item>
      <p v-if="log" class="detail-meta">
        创建 {{ log.created_at }} · 更新 {{ log.last_updated_at || log.updated_at }}
      </p>
      <p v-if="log?.location" class="detail-meta">地点 {{ log.location }}</p>
    </n-form>
    <template #action>
      <n-button tertiary type="error" :loading="loading" :disabled="!log" @click="remove">删除</n-button>
      <n-button @click="emit('update:show', false)">取消</n-button>
      <n-button type="primary" :loading="loading" :disabled="!canSave" @click="save">保存</n-button>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { LifeLogInfo } from '@/types/api'

const props = defineProps<{ show: boolean; log?: LifeLogInfo | null; loading?: boolean }>()
const emit = defineEmits<{
  'update:show': [value: boolean]
  save: [payload: { id: number; content: string; occurred_at: string }]
  delete: [id: number]
}>()

const form = ref({ content: '', occurred_at: '' })
const canSave = computed(() => Boolean(props.log && form.value.content.trim()))

watch(
  () => [props.log, props.show] as const,
  ([log, show]) => {
    if (!log || !show) return
    form.value = {
      content: log.content,
      occurred_at: log.occurred_at
    }
  },
  { immediate: true }
)

function save() {
  if (!props.log || !canSave.value) return
  emit('save', {
    id: props.log.id,
    content: form.value.content.trim(),
    occurred_at: form.value.occurred_at
  })
}

function remove() {
  if (!props.log) return
  emit('delete', props.log.id)
}
</script>
