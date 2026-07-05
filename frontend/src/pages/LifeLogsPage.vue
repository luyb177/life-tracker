<template>
  <div class="page">
    <PageHeader title="生活记录" description="按日期查看、补充和管理每天发生的事情。">
      <n-date-picker
        v-model:formatted-value="date"
        value-format="yyyy-MM-dd"
        type="date"
        :is-date-disabled="isFutureDateTimestamp"
      />
    </PageHeader>

    <section class="panel">
      <div class="section-title-row">
        <div>
          <p class="eyebrow">Logs</p>
          <h2>{{ date }} 的记录</h2>
        </div>
        <n-button secondary :loading="loading" @click="load">刷新</n-button>
      </div>

      <EmptyState
        v-if="logs.length === 0"
        title="这天还没有生活记录"
        description="回到今天页面可以快速补一条。"
        :icon="BookOpenText"
      />
      <div v-else class="record-list life-log-record-list">
        <article v-for="log in logs" :key="log.id" class="record-row interactive-record" @click="openDetail(log)">
          <time>{{ log.occurred_at }}</time>
          <p>{{ log.content }}</p>
          <div class="tag-row compact">
            <span v-for="tag in log.tags" :key="tag.id" class="chip readonly">{{ tag.name }}</span>
          </div>
        </article>
      </div>
    </section>

    <LifeLogDetailModal
      v-model:show="showDetailModal"
      :log="selectedLog"
      :loading="saving"
      @save="saveSelected"
      @delete="deleteSelected"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useMessage } from 'naive-ui'
import { BookOpenText } from '@lucide/vue'
import PageHeader from '@/components/common/PageHeader.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import LifeLogDetailModal from '@/components/records/LifeLogDetailModal.vue'
import { deleteLifeLog, getLifeLogsByDate, updateLifeLog } from '@/api/lifeLog'
import type { LifeLogInfo } from '@/types/api'
import { formatDate, isFutureDateTimestamp } from '@/utils/date'

const message = useMessage()
const date = ref(formatDate())
const logs = ref<LifeLogInfo[]>([])
const loading = ref(false)
const saving = ref(false)
const showDetailModal = ref(false)
const selectedLog = ref<LifeLogInfo | null>(null)

function openDetail(log: LifeLogInfo) {
  selectedLog.value = log
  showDetailModal.value = true
}

async function saveSelected(payload: { id: number; content: string; occurred_at: string }) {
  saving.value = true
  try {
    await updateLifeLog(payload)
    message.success('生活记录已更新')
    showDetailModal.value = false
    await load()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '更新失败')
  } finally {
    saving.value = false
  }
}

async function deleteSelected(id: number) {
  if (!window.confirm('确认删除这条生活记录吗？')) return
  saving.value = true
  try {
    await deleteLifeLog(id)
    message.success('生活记录已删除')
    showDetailModal.value = false
    await load()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '删除失败')
  } finally {
    saving.value = false
  }
}

async function load() {
  loading.value = true
  try {
    logs.value = (await getLifeLogsByDate(date.value)).list
  } catch (error) {
    message.error(error instanceof Error ? error.message : '加载记录失败')
  } finally {
    loading.value = false
  }
}

watch(date, load)
onMounted(load)
</script>
