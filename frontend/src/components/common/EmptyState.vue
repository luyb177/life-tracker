<template>
  <div class="empty-state">
    <component :is="icon" v-if="icon" :size="24" />
    <strong>{{ title }}</strong>
    <p>{{ description }}</p>
  </div>
</template>

<script setup lang="ts">
import type { Component } from 'vue'

defineProps<{
  title: string
  description: string
  icon?: Component
}>()
</script>

