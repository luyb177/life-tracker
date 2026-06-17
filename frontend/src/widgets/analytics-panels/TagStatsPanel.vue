<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { NTag } from 'naive-ui'
import EmptyState from '@/shared/ui/empty-state/EmptyState.vue'
import ChartCard from '@/shared/ui/chart-card/ChartCard.vue'
import { getTagStats } from '@/entities/summary/api/summary.api'
import { summaryKeys } from '@/entities/summary/api/summary.keys'

const props = defineProps<{ start: string; end: string }>()

const { data, isLoading } = useQuery({
  queryKey: computed(() => summaryKeys.tagStats({ start: props.start, end: props.end })),
  queryFn: () => getTagStats({ start: props.start, end: props.end }),
  enabled: computed(() => !!props.start),
})

const tags = computed(() => (data.value as any)?.tags || [])
</script>
<template>
  <ChartCard title="标签频次" :loading="isLoading" :is-empty="!tags.length">
    <div v-if="tags.length" style="display: flex; gap: 8px; flex-wrap: wrap">
      <NTag v-for="t in tags" :key="t.tag" type="success">{{ t.tag }} ({{ t.count }})</NTag>
    </div>
    <EmptyState v-else title="暂无标签" />
  </ChartCard>
</template>
