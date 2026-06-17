<script setup lang="ts">
import { NTag } from 'naive-ui'
import SummaryCard from '@/shared/ui/summary-card/SummaryCard.vue'
import type { SummaryInfo } from '@/entities/summary/api/summary.api'
import { PERIOD_LABELS } from '@/shared/constants'

defineProps<{ items: SummaryInfo[]; periodType: number }>()
const emit = defineEmits<{ edit: [id: number]; delete: [id: number] }>()
</script>
<template>
  <div v-if="items.length">
    <h4 style="margin: 16px 0 8px"><NTag type="success" size="small">我的{{ PERIOD_LABELS[periodType] }}</NTag></h4>
    <SummaryCard v-for="s in items" :key="s.id" :item="s" show-actions @edit="emit('edit', s.id)" @delete="emit('delete', s.id)" />
  </div>
</template>
