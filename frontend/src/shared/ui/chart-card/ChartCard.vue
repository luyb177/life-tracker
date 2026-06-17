<script setup lang="ts">
import { NSpin } from 'naive-ui'
import EmptyState from '@/shared/ui/empty-state/EmptyState.vue'
defineProps<{ title: string; loading?: boolean; isEmpty?: boolean }>()
</script>
<template>
  <div style="margin-bottom: 24px">
    <h4 style="margin: 0 0 8px">{{ title }}</h4>
    <NSpin :show="loading">
      <template v-if="!isEmpty"><slot /></template>
      <EmptyState v-else title="暂无数据" />
    </NSpin>
  </div>
</template>
