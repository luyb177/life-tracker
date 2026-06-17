<script setup lang="ts">
import { NTag } from 'naive-ui'
import SummaryCard from '@/shared/ui/summary-card/SummaryCard.vue'
import type { SummaryInfo } from '@/entities/summary/api/summary.api'
import { PERIOD_LABELS } from '@/shared/constants'

defineProps<{ items: SummaryInfo[]; periodType: number }>()
</script>
<template>
  <div v-if="items.length">
    <h4 style="margin-bottom: 8px"><NTag type="info" size="small">AI {{ PERIOD_LABELS[periodType] }}</NTag></h4>
    <SummaryCard v-for="s in items" :key="s.id" :item="s" />
  </div>
</template>
