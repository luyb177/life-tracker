<script setup lang="ts">
import { NButton } from 'naive-ui'
import FilterBar from '@/shared/ui/filter-bar/FilterBar.vue'
import PageHeader from '@/shared/ui/page-header/PageHeader.vue'
import PageSection from '@/shared/ui/page-section/PageSection.vue'
import SummaryEditorModal from '@/shared/ui/summary-editor-modal/SummaryEditorModal.vue'
import EmptyState from '@/shared/ui/empty-state/EmptyState.vue'
import AISummaryList from '@/widgets/summary-lists/AISummaryList.vue'
import UserSummaryList from '@/widgets/summary-lists/UserSummaryList.vue'
import { PERIOD_LABELS } from '@/shared/constants'
import { useSummaries } from '@/features/summary-create/useSummaries'

const { periodType, showEditor, editTarget, submitting, isEditing, editorTitle, isLoading, list, aiSummaries, userSummaries, openCreate, openEdit, closeEditor, handleCreate, handleUpdate, handleDelete, handleGenerate } = useSummaries()
</script>
<template>
  <PageHeader title="总结" description="AI 周期总结与个人复盘" />
  <PageSection title="总结操作" description="先选择周期，再决定是自己写总结，还是生成 AI 总结。">
    <FilterBar v-model:period-type="periodType" show-period>
      <NButton size="small" @click="openCreate">写{{ PERIOD_LABELS[periodType] }}</NButton>
      <NButton size="small" type="primary" @click="handleGenerate">AI 生成</NButton>
    </FilterBar>
  </PageSection>

  <PageSection title="总结结果" description="把 AI 总结和自己的总结放在一起看，方便横向对比与复盘。">
    <div v-if="isLoading">加载中...</div>
    <template v-else-if="list.length">
      <AISummaryList :items="aiSummaries" :period-type="periodType" />
      <UserSummaryList :items="userSummaries" :period-type="periodType" @edit="openEdit" @delete="handleDelete" />
    </template>
    <EmptyState v-else :title="`暂无${PERIOD_LABELS[periodType]}`" />
  </PageSection>

  <SummaryEditorModal
    v-model:show="showEditor"
    :title="editorTitle"
    :item="editTarget"
    :submitting="submitting"
    :show-period-select="true"
    :lock-period-fields="isEditing"
    :initial-period-type="periodType"
    @update:show="!$event && closeEditor()"
    @submit="isEditing ? handleUpdate($event) : handleCreate($event)"
  />
</template>
