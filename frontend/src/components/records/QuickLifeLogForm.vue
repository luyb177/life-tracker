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
      <button v-if="!addingTag" type="button" class="chip add-chip" @click="addingTag = true">
        + 标签
      </button>
      <n-input
        v-else
        ref="tagInputRef"
        v-model:value="customTag"
        class="inline-tag-input"
        size="small"
        placeholder="输入标签"
        @keyup.enter="addCustomTag"
        @blur="addCustomTag"
      />
    </div>

    <n-button type="primary" block :loading="loading" :disabled="!content.trim()" @click="submit">
      保存生活记录
    </n-button>
  </n-form>
</template>

<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'
import { useMessage } from 'naive-ui'
import { createLifeLog } from '@/api/lifeLog'
import { formatDate, formatDateTime } from '@/utils/date'

const props = defineProps<{ date: string }>()
const emit = defineEmits<{ created: [] }>()

const message = useMessage()
const content = ref('')
const loading = ref(false)
const timeValue = ref(Date.now())
const defaultTags = ['工作', '学习', '运动', '娱乐', '睡眠']
const presetTags = ref(readTags())
const selectedTags = ref<string[]>([])
const addingTag = ref(false)
const customTag = ref('')
const tagInputRef = ref<{ focus: () => void } | null>(null)

function readTags() {
  const raw = localStorage.getItem('life-tracker.life-tags')
  if (!raw) return defaultTags
  try {
    const saved = JSON.parse(raw) as string[]
    return Array.from(new Set([...defaultTags, ...saved])).filter(Boolean)
  } catch {
    return defaultTags
  }
}

function toggleTag(tag: string) {
  selectedTags.value = selectedTags.value.includes(tag)
    ? selectedTags.value.filter((item) => item !== tag)
    : [...selectedTags.value, tag]
}

async function addCustomTag() {
  const tag = customTag.value.trim()
  if (!tag) {
    addingTag.value = false
    return
  }
  if (!presetTags.value.includes(tag)) {
    presetTags.value = [...presetTags.value, tag]
    localStorage.setItem('life-tracker.life-tags', JSON.stringify(presetTags.value))
  }
  if (!selectedTags.value.includes(tag)) {
    selectedTags.value = [...selectedTags.value, tag]
  }
  customTag.value = ''
  addingTag.value = false
}

function buildOccurredAt() {
  const date = props.date || formatDate()
  const time = new Date(timeValue.value)
  return formatDateTime(new Date(`${date}T${time.toTimeString().slice(0, 8)}`))
}

watch(addingTag, async (value) => {
  if (!value) return
  await nextTick()
  tagInputRef.value?.focus()
})

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
