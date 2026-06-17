<script setup lang="ts">
import StatCard from '@/shared/ui/stat-card/StatCard.vue'

defineProps<{ total: number; count: number; loading?: boolean }>()
</script>
<template>
  <div style="display: flex; gap: 12px; margin-bottom: 16px">
    <StatCard title="今日支出" :value="total.toFixed(2)" unit="元" :loading="loading" />
    <StatCard title="笔数" :value="count" unit="笔" :loading="loading" />
  </div>
</template>
