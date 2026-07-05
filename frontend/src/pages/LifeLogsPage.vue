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

    <n-modal v-model:show="showDetailModal" preset="dialog" title="生活记录详情" class="detail-modal">
      <n-form class="detail-form" :show-label="true" label-placement="top">
        <n-form-item label="发生时间">
          <n-date-picker
            v-model:formatted-value="editForm.occurred_at"
            type="datetime"
            value-format="yyyy-MM-dd HH:mm:ss"
            class="full-control"
          />
        </n-form-item>
        <n-form-item label="内容">
          <n-input v-model:value="editForm.content" type="textarea" :autosize="{ minRows: 5, maxRows: 9 }" />
        </n-form-item>
        <p v-if="selectedLog" class="detail-meta">
          创建 {{ selectedLog.created_at }} · 更新 {{ selectedLog.last_updated_at || selectedLog.updated_at }}
        </p>
      </n-form>
      <template #action>
        <n-button tertiary type="error" :loading="saving" @click="deleteSelected">删除</n-button>
        <n-button @click="showDetailModal = false">取消</n-button>
        <n-button type="primary" :loading="saving" :disabled="!editForm.content.trim()" @click="saveSelected">
          保存
        </n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useMessage } from 'naive-ui'
import { BookOpenText } from '@lucide/vue'
import PageHeader from '@/components/common/PageHeader.vue'
import EmptyState from '@/components/common/EmptyState.vue'
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
const editForm = ref({ content: '', occurred_at: '' })

function openDetail(log: LifeLogInfo) {
  selectedLog.value = log
  editForm.value = {
    content: log.content,
    occurred_at: log.occurred_at
  }
  showDetailModal.value = true
}

async function saveSelected() {
  if (!selectedLog.value || !editForm.value.content.trim()) return
  saving.value = true
  try {
    await updateLifeLog({
      id: selectedLog.value.id,
      content: editForm.value.content.trim(),
      occurred_at: editForm.value.occurred_at
    })
    message.success('生活记录已更新')
    showDetailModal.value = false
    await load()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '更新失败')
  } finally {
    saving.value = false
  }
}

async function deleteSelected() {
  if (!selectedLog.value || !window.confirm('确认删除这条生活记录吗？')) return
  saving.value = true
  try {
    await deleteLifeLog(selectedLog.value.id)
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
