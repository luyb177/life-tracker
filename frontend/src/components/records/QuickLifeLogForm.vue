<template>
  <n-form class="quick-form" :show-label="false" @submit.prevent="submit">
    <div class="section-title-row">
      <div>
        <p class="eyebrow">生活</p>
        <h2>快速记录</h2>
      </div>
      <n-time-picker v-model:value="timeValue" size="small" format="HH:mm" />
    </div>

    <n-input
      v-model:value="content"
      type="textarea"
      placeholder="刚刚做了什么？"
      :autosize="{ minRows: 4, maxRows: 6 }"
    />

    <div class="tag-row">
      <button
        v-for="tag in presetTags"
        :key="tag"
        type="button"
        class="chip"
        :class="{ active: selectedTags.includes(tag) }"
        @click="toggleTag(tag)"
      >
        {{ tag }}
      </button>
    </div>

    <n-button type="primary" block :loading="loading" :disabled="!content.trim()" @click="submit">
      保存生活记录
    </n-button>
  </n-form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { createLifeLog } from '@/api/lifeLog'
import { formatDate, formatDateTime } from '@/utils/date'

const props = defineProps<{ date: string }>()
const emit = defineEmits<{ created: [] }>()

const message = useMessage()
const content = ref('')
const loading = ref(false)
const timeValue = ref(Date.now())
const presetTags = ['工作', '学习', '运动', '娱乐', '睡眠']
const selectedTags = ref<string[]>([])

function toggleTag(tag: string) {
  selectedTags.value = selectedTags.value.includes(tag)
    ? selectedTags.value.filter((item) => item !== tag)
    : [...selectedTags.value, tag]
}

function buildOccurredAt() {
  const date = props.date || formatDate()
  const time = new Date(timeValue.value)
  return formatDateTime(new Date(`${date}T${time.toTimeString().slice(0, 8)}`))
}

async function submit() {
  if (!content.value.trim()) return
  loading.value = true
  try {
    await createLifeLog({
      content: content.value.trim(),
      occurred_at: buildOccurredAt(),
      tags: selectedTags.value.map((name) => ({ id: 0, name }))
    })
    content.value = ''
    selectedTags.value = []
    message.success('生活记录已保存')
    emit('created')
  } finally {
    loading.value = false
  }
}
</script>

