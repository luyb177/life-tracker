<script setup lang="ts">
defineProps<{
  title: string
  description?: string
}>()
</script>
<template>
  <section style="margin-top: 24px">
    <div style="margin-bottom: 12px">
      <h3 style="margin: 0; font-size: 17px">{{ title }}</h3>
      <p v-if="description" style="margin: 4px 0 0; font-size: 13px; color: var(--n-text-color-3)">{{ description }}</p>
    </div>
    <slot />
  </section>
</template>
