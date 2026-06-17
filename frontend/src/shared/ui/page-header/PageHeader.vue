<script setup lang="ts">
defineProps<{ title: string; description?: string }>()
</script>
<template>
  <div style="margin-bottom: 16px">
    <h2 style="margin: 0; font-size: 20px">{{ title }}</h2>
    <p v-if="description" style="margin: 4px 0 0; font-size: 13px; color: var(--n-text-color-3)">{{ description }}</p>
  </div>
</template>
