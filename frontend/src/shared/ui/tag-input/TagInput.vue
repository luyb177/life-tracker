<script setup lang="ts">
import { computed } from 'vue'
import { NInput, NTag } from 'naive-ui'
import { parseTags } from '@/shared/utils/format'

const props = defineProps<{
  modelValue: string
  placeholder?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const tags = computed(() => parseTags(props.modelValue))
</script>
<template>
  <div>
    <NInput
      :value="modelValue"
      :placeholder="placeholder || '标签，逗号分隔'"
      @update:value="emit('update:modelValue', $event)"
    />
    <div v-if="tags.length" style="display: flex; gap: 6px; flex-wrap: wrap; margin-top: 8px">
      <NTag v-for="tag in tags" :key="tag" size="tiny" :bordered="false">{{ tag }}</NTag>
    </div>
  </div>
</template>
