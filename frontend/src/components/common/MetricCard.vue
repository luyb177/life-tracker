<template>
  <section class="metric-card">
    <div class="metric-icon" :class="tone">
      <component :is="icon" :size="20" />
    </div>
    <div>
      <p>{{ label }}</p>
      <strong>{{ value }}</strong>
      <span v-if="hint">{{ hint }}</span>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { Component } from 'vue'

defineProps<{
  label: string
  value: string
  hint?: string
  icon: Component
  tone?: 'teal' | 'coral' | 'amber' | 'green'
}>()
</script>

