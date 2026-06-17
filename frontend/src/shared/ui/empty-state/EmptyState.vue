<script setup lang="ts">
defineProps<{ title?: string; description?: string }>()
</script>
<template>
  <div style="text-align: center; padding: 48px 16px; color: var(--n-text-color-3)">
    <div style="font-size: 48px; margin-bottom: 12px">📭</div>
    <div style="font-size: 16px; font-weight: 500">{{ title || '暂无数据' }}</div>
    <div v-if="description" style="font-size: 13px; margin-top: 4px">{{ description }}</div>
  </div>
</template>
