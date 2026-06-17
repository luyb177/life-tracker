<script setup lang="ts">
import { NButton } from 'naive-ui'
import FilterBar from '@/shared/ui/filter-bar/FilterBar.vue'
import PageHeader from '@/shared/ui/page-header/PageHeader.vue'
import PageSection from '@/shared/ui/page-section/PageSection.vue'
import SummaryEditorModal from '@/shared/ui/summary-editor-modal/SummaryEditorModal.vue'
import EmptyState from '@/shared/ui/empty-state/EmptyState.vue'
import AISummaryList from '@/widgets/summary-lists/AISummaryList.vue'
import UserSummaryList from '@/widgets/summary-lists/UserSummaryList.vue'
import { useRecords } from '@/features/record-create/useRecords'
import { PERIOD_TYPE } from '@/shared/constants'

const { date, dateStr, showEditor, editTarget, submitting, isLoading, list, userRecords, aiSummaries, editorTitle, openCreate, openEdit, closeEditor, handleSubmit, handleDelete } = useRecords()
</script>
<template>
  <PageHeader title="生活记录" description="记录每天做了什么" />
  <PageSection title="记录操作" description="先切换日期，再补写或整理当天的生活记录。">
    <FilterBar v-model:date="date" show-date>
      <NButton type="primary" size="small" @click="openCreate">写记录</NButton>
    </FilterBar>
  </PageSection>

  <PageSection title="当天记录" :description="`${dateStr} 的个人记录与 AI 日报会一起展示，方便回看当天内容。`">
    <div v-if="isLoading">加载中...</div>
    <template v-else-if="list.length">
      <UserSummaryList :items="userRecords" :period-type="PERIOD_TYPE.DAY" @edit="openEdit" @delete="handleDelete" />
      <AISummaryList :items="aiSummaries" :period-type="PERIOD_TYPE.DAY" />
    </template>
    <EmptyState v-else title="当天暂无记录" description="点击「写记录」开始记录生活" />
  </PageSection>

  <SummaryEditorModal
    v-model:show="showEditor"
    :title="editorTitle"
    :item="editTarget"
    :submitting="submitting"
    :initial-period-type="PERIOD_TYPE.DAY"
    :initial-period-start="dateStr"
    @update:show="!$event && closeEditor()"
    @submit="handleSubmit"
  />
</template>
