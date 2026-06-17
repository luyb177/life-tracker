<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { NTag } from 'naive-ui'
import EmptyState from '@/shared/ui/empty-state/EmptyState.vue'
import { getTagTrend } from '@/entities/summary/api/summary.api'
import { summaryKeys } from '@/entities/summary/api/summary.keys'

const props = defineProps<{ start: string; end: string }>()

const { data } = useQuery({
  queryKey: computed(() => summaryKeys.tagTrend({ start: props.start, end: props.end })),
  queryFn: () => getTagTrend({ start: props.start, end: props.end }),
  enabled: computed(() => !!props.start),
})
const months = computed(() => (data.value as any)?.months || [])
</script>
<template>
  <div style="margin-top: 16px">
    <h4 style="margin-bottom: 8px">标签按月趋势</h4>
    <div v-if="months.length">
      <div v-for="m in months" :key="m.month" style="margin-bottom: 12px">
        <h5 style="margin: 0 0 4px">{{ m.month }}</h5>
        <div style="display: flex; gap: 6px; flex-wrap: wrap">
          <NTag v-for="t in m.tags.slice(0, 8)" :key="t.tag" size="tiny">{{ t.tag }}({{ t.count }})</NTag>
        </div>
      </div>
    </div>
    <EmptyState v-else title="暂无趋势数据" />
  </div>
</template>
