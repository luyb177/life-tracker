<script setup lang="ts">
import { NButton, NModal } from 'naive-ui'
import FilterBar from '@/shared/ui/filter-bar/FilterBar.vue'
import PageHeader from '@/shared/ui/page-header/PageHeader.vue'
import PageSection from '@/shared/ui/page-section/PageSection.vue'
import ExpenseForm from '@/shared/ui/expense-form/ExpenseForm.vue'
import { useExpenses } from '@/features/expense-create/useExpenses'
import ExpenseStatsPanel from '@/widgets/expense-day-panel/ExpenseStatsPanel.vue'
import ExpenseListPanel from '@/widgets/expense-day-panel/ExpenseListPanel.vue'
import ExpenseCategoryManager from '@/widgets/expense-day-panel/ExpenseCategoryManager.vue'

const { date, showModal, editTarget, submitting, showCatModal, newCatName, isLoading, categories, list, total, expenseModalTitle, openCreate, openEdit, closeEditor, openCategoryManager, closeCategoryManager, handleSubmit, handleDelete, handleCreateCategory, handleDeleteCategory } = useExpenses()
</script>
<template>
  <PageHeader title="支出" description="记录每日开销" />
  <PageSection title="记录操作" description="按天记录支出，并在需要时维护自己的分类。">
    <FilterBar v-model:date="date" show-date>
      <NButton type="primary" size="small" @click="openCreate">记一笔</NButton>
      <NButton size="small" @click="openCategoryManager">分类</NButton>
    </FilterBar>
  </PageSection>

  <PageSection title="当日概览" description="先看当天总支出和支出笔数，再决定是否继续查看明细。">
    <ExpenseStatsPanel :total="total" :count="list.length" :loading="isLoading" />
  </PageSection>

  <PageSection title="支出明细" description="这里保留当天的具体消费明细，方便回顾每一笔开销。">
    <div v-if="isLoading">加载中...</div>
    <ExpenseListPanel v-else :list="list" :total="total" @edit="openEdit" @delete="handleDelete" />
  </PageSection>

  <NModal v-model:show="showModal" :title="expenseModalTitle">
    <div style="padding: 16px"><ExpenseForm :edit-id="editTarget?.id" :initial-category-id="editTarget?.category?.id" :initial-amount="editTarget?.amount" :initial-note="editTarget?.note" :submitting="submitting" @submit="handleSubmit" @cancel="closeEditor" /></div>
  </NModal>

  <NModal v-model:show="showCatModal" title="管理分类">
    <ExpenseCategoryManager
      :categories="categories"
      :new-category-name="newCatName"
      @update:new-category-name="newCatName = $event"
      @create="handleCreateCategory"
      @delete="handleDeleteCategory"
    />
  </NModal>
</template>
