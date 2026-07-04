<template>
  <div class="page">
    <PageHeader title="支出" description="查看每日支出、分类结构和近期趋势。">
      <n-date-picker v-model:formatted-value="date" value-format="yyyy-MM-dd" type="date" />
    </PageHeader>

    <div class="metric-grid">
      <MetricCard label="当日合计" :value="formatYuan(total)" hint="已退款不计入" :icon="Wallet" tone="coral" />
      <MetricCard label="支出笔数" :value="`${expenses.length} 笔`" hint="选中日期" :icon="ReceiptText" tone="teal" />
      <MetricCard label="分类数量" :value="`${categories.length} 个`" hint="系统与自定义" :icon="Tags" tone="amber" />
    </div>

    <div class="today-grid wide-left">
      <section class="panel">
        <div class="section-title-row">
          <div>
            <p class="eyebrow">Records</p>
            <h2>{{ date }} 的支出</h2>
          </div>
          <n-button secondary :loading="loading" @click="load">刷新</n-button>
        </div>

        <EmptyState
          v-if="expenses.length === 0"
          title="这天还没有支出"
          description="回到今天页面可以快速记一笔花销。"
          :icon="ReceiptText"
        />
        <div v-else class="record-list">
          <article v-for="item in expenses" :key="item.id" class="record-row split">
            <div>
              <time>{{ item.occurred_at }}</time>
              <p>{{ item.category.name }} · {{ item.note || '无备注' }}</p>
            </div>
            <strong>{{ formatYuan(item.amount) }}</strong>
          </article>
        </div>
      </section>

      <section class="panel">
        <div class="section-title-row">
          <div>
            <p class="eyebrow">Categories</p>
            <h2>支出分类</h2>
          </div>
        </div>
        <div class="category-list">
          <span v-for="category in categories" :key="category.id" class="chip readonly">
            {{ category.name }}
          </span>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useMessage } from 'naive-ui'
import { ReceiptText, Tags, Wallet } from '@lucide/vue'
import PageHeader from '@/components/common/PageHeader.vue'
import MetricCard from '@/components/common/MetricCard.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import { getExpensesByDate, listExpenseCategories } from '@/api/expense'
import type { ExpenseCategoryInfo, ExpenseLogInfo } from '@/types/api'
import { formatDate } from '@/utils/date'
import { formatYuan } from '@/utils/money'

const message = useMessage()
const date = ref(formatDate())
const expenses = ref<ExpenseLogInfo[]>([])
const categories = ref<ExpenseCategoryInfo[]>([])
const total = ref(0)
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const [expenseResp, categoryResp] = await Promise.all([
      getExpensesByDate(date.value),
      listExpenseCategories()
    ])
    expenses.value = expenseResp.list
    total.value = expenseResp.total
    categories.value = categoryResp.categories
  } catch (error) {
    message.error(error instanceof Error ? error.message : '加载支出失败')
  } finally {
    loading.value = false
  }
}

watch(date, load)
onMounted(load)
</script>

