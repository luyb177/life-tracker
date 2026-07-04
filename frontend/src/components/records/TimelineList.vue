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

    <div v-else class="timeline today-timeline-list">
      <article v-for="item in items" :key="item.id" class="timeline-item" :class="[item.type, { refunded: item.refunded }]">
        <div class="timeline-time">{{ item.time }}</div>
        <div>
          <div class="timeline-title-row">
            <strong>{{ item.title }}</strong>
            <div class="timeline-side">
              <span v-if="item.amount" class="timeline-amount">{{ item.amount }}</span>
              <span v-if="item.refunded" class="status-badge">已退款</span>
              <n-button
                v-else-if="item.canRefund"
                size="tiny"
                tertiary
                type="error"
                @click="$emit('refund', item.sequence)"
              >
                退款
              </n-button>
            </div>
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
  sortAt: string
  createdAt: string
  sequence: number
  title: string
  description: string
  amount?: string
  canRefund?: boolean
  refunded?: boolean
  tags?: string[]
}

defineProps<{ items: TimelineItem[] }>()
defineEmits<{ refund: [id: number] }>()
</script>
