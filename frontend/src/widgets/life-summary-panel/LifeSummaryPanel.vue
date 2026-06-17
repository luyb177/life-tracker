<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { NButton } from 'naive-ui'
import SummaryCard from '@/shared/ui/summary-card/SummaryCard.vue'
import EmptyState from '@/shared/ui/empty-state/EmptyState.vue'
import { getSummaryList, generateAISummary } from '@/entities/summary/api/summary.api'
import { summaryKeys } from '@/entities/summary/api/summary.keys'
import { PERIOD_TYPE } from '@/shared/constants'

const { data } = useQuery({
  queryKey: summaryKeys.list({ page_size: 5, period_type: PERIOD_TYPE.LIFE }),
  queryFn: () => getSummaryList({ page_size: 5, period_type: PERIOD_TYPE.LIFE }),
})
const list = computed(() => (data.value as any)?.list || [])

async function handleGenerate() {
  await generateAISummary(PERIOD_TYPE.LIFE, '2000-01-01')
}
</script>
<template>
  <div>
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px">
      <h4 style="margin: 0">📊 人生总结</h4>
      <NButton size="tiny" @click="handleGenerate">生成</NButton>
    </div>
    <template v-if="list.length">
      <SummaryCard v-for="s in list" :key="s.id" :item="s" />
    </template>
    <EmptyState v-else title="暂无人生总结" />
  </div>
</template>
