<script setup lang="ts">
import { NButton } from 'naive-ui'
import StatCard from '@/shared/ui/stat-card/StatCard.vue'

defineProps<{
  recordCount: number
  totalExpense: number
  expenseCount: number
}>()

const emit = defineEmits<{
  records: []
  expenses: []
  summaries: []
}>()
</script>
<template>
  <div style="margin-bottom: 20px">
    <div style="display: flex; gap: 12px; flex-wrap: wrap; margin-bottom: 20px">
      <StatCard title="今日记录" :value="recordCount" unit="条" />
      <StatCard title="今日支出" :value="totalExpense.toFixed(2)" unit="元" />
      <StatCard title="支出笔数" :value="expenseCount" unit="笔" />
    </div>
    <div style="display: flex; gap: 8px; flex-wrap: wrap">
      <NButton type="primary" @click="emit('records')">写记录</NButton>
      <NButton @click="emit('expenses')">记支出</NButton>
      <NButton @click="emit('summaries')">看总结</NButton>
    </div>
  </div>
</template>
