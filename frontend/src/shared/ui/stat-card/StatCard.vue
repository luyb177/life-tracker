<script setup lang="ts">
defineProps<{ title: string; value: string | number; unit?: string; loading?: boolean }>()
</script>
<template>
  <div style="background: var(--n-color-target); border-radius: 8px; padding: 16px; min-width: 140px">
    <div style="font-size: 13px; color: var(--n-text-color-3); margin-bottom: 8px">{{ title }}</div>
    <div style="font-size: 28px; font-weight: 700">
      <template v-if="loading">...</template>
      <template v-else>{{ value }}<span v-if="unit" style="font-size: 14px; font-weight: 400"> {{ unit }}</span></template>
    </div>
  </div>
</template>
