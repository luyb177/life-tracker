<script setup lang="ts">
import TodaySummaryPanel from '@/widgets/today-summary-panel/TodaySummaryPanel.vue'
import TodayExpensePanel from '@/widgets/today-expense-panel/TodayExpensePanel.vue'

defineProps<{ date: string }>()
</script>
<template>
  <div style="margin-top: 20px; display: grid; gap: 16px">
    <TodaySummaryPanel :date="date" />
    <TodayExpensePanel :date="date" />
  </div>
</template>
