<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { NSpin, NTag } from 'naive-ui'
import SummaryCard from '@/shared/ui/summary-card/SummaryCard.vue'
import EmptyState from '@/shared/ui/empty-state/EmptyState.vue'
import { getDaySummaries, SummaryInfo } from '@/entities/summary/api/summary.api'
import { summaryKeys } from '@/entities/summary/api/summary.keys'
import { isAISummary, isUserRecord } from '@/shared/utils/format'

const props = defineProps<{ date: string }>()

const { data, isLoading } = useQuery({
  queryKey: computed(() => summaryKeys.day(props.date)),
  queryFn: () => getDaySummaries(props.date),
  enabled: computed(() => !!props.date),
})

const list = computed(() => ((data.value as any)?.list || []) as SummaryInfo[])
const userRecords = computed(() => list.value.filter(s => isUserRecord(s.source)))
const aiSummaries = computed(() => list.value.filter(s => isAISummary(s.source)))
</script>
<template>
  <NSpin :show="isLoading">
    <template v-if="userRecords.length">
      <h4 style="display: flex; align-items: center; gap: 8px; margin-bottom: 8px">
        <NTag type="success" size="small">我的记录</NTag>
        <span style="font-size: 13px; color: var(--n-text-color-3)">{{ userRecords.length }} 条</span>
      </h4>
      <SummaryCard v-for="s in userRecords" :key="s.id" :item="s" />
    </template>
    <template v-if="aiSummaries.length">
      <h4 style="margin: 16px 0 8px"><NTag type="info" size="small">AI 日报</NTag></h4>
      <SummaryCard v-for="s in aiSummaries" :key="s.id" :item="s" />
    </template>
    <EmptyState v-if="!list.length" title="今日暂无记录" description="前往「生活记录」页开始记录" />
  </NSpin>
</template>
