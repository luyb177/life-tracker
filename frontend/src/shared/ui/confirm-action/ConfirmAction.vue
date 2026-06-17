<script setup lang="ts">
import { NPopconfirm, NButton } from 'naive-ui'
defineProps<{ text?: string; confirmText?: string }>()
defineEmits<{ confirm: [] }>()
</script>
<template>
  <NPopconfirm @positive-click="$emit('confirm')">
    <template #trigger><NButton text size="tiny" type="error">{{ text || '删除' }}</NButton></template>
    {{ confirmText || '确定删除？' }}
  </NPopconfirm>
</template>
