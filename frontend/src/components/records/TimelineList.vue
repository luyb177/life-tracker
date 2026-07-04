<template>
  <section class="panel">
    <div class="section-title-row">
      <div>
        <p class="eyebrow">今天</p>
        <h2>时间线</h2>
      </div>
      <span class="muted">{{ items.length }} 条</span>
    </div>

    <EmptyState
      v-if="items.length === 0"
      title="还没有记录"
      description="先写下一条生活记录或支出，今天就有了清晰的起点。"
      :icon="Clock3"
    />

    <div v-else class="timeline">
      <article v-for="item in items" :key="item.id" class="timeline-item" :class="item.type">
        <div class="timeline-time">{{ item.time }}</div>
        <div>
          <div class="timeline-title-row">
            <strong>{{ item.title }}</strong>
            <span v-if="item.amount" class="timeline-amount">{{ item.amount }}</span>
          </div>
          <p>{{ item.description }}</p>
          <div v-if="item.tags?.length" class="tag-row compact">
            <span v-for="tag in item.tags" :key="tag" class="chip readonly">{{ tag }}</span>
          </div>
        </div>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
import { Clock3 } from '@lucide/vue'
import EmptyState from '@/components/common/EmptyState.vue'

export interface TimelineItem {
  id: string
  type: 'life' | 'expense'
  time: string
  title: string
  description: string
  amount?: string
  tags?: string[]
}

defineProps<{ items: TimelineItem[] }>()
</script>

