<template>
  <n-form class="quick-form" :show-label="false" @submit.prevent="submit">
    <div class="section-title-row">
      <div>
        <p class="eyebrow">支出</p>
        <h2>记一笔花销</h2>
      </div>
      <n-time-picker v-model:value="timeValue" size="small" format="HH:mm" />
    </div>

    <div class="expense-grid">
      <n-input-number
        v-model:value="amount"
        placeholder="金额"
        :precision="2"
        :min="0"
        class="amount-input"
      >
        <template #prefix>¥</template>
      </n-input-number>
      <n-select
        v-model:value="categoryId"
        :options="categoryOptions"
        placeholder="分类"
        :loading="loadingCategories"
      />
    </div>

    <n-input v-model:value="note" placeholder="备注，可选" />

    <n-button type="primary" block :loading="loading" :disabled="!canSubmit" @click="submit">
      保存支出
    </n-button>
  </n-form>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useMessage } from 'naive-ui'
import { createExpense, listExpenseCategories } from '@/api/expense'
import { formatDate, formatDateTime } from '@/utils/date'
import { yuanToFen } from '@/utils/money'
import type { ExpenseCategoryInfo } from '@/types/api'

const props = defineProps<{ date: string }>()
const emit = defineEmits<{ created: [] }>()

const message = useMessage()
const amount = ref<number | null>(null)
const categoryId = ref<number | null>(null)
const note = ref('')
const timeValue = ref(Date.now())
const loading = ref(false)
const loadingCategories = ref(false)
const categories = ref<ExpenseCategoryInfo[]>([])

const categoryOptions = computed(() =>
  categories.value.map((item) => ({ label: item.name, value: item.id }))
)
const canSubmit = computed(() => Boolean(amount.value && categoryId.value))

function buildOccurredAt() {
  const date = props.date || formatDate()
  const time = new Date(timeValue.value)
  return formatDateTime(new Date(`${date}T${time.toTimeString().slice(0, 8)}`))
}

async function loadCategories() {
  loadingCategories.value = true
  try {
    const resp = await listExpenseCategories()
    categories.value = resp.categories
    categoryId.value = resp.categories[0]?.id ?? null
  } finally {
    loadingCategories.value = false
  }
}

async function submit() {
  if (!canSubmit.value || !categoryId.value) return
  loading.value = true
  try {
    await createExpense({
      category_id: categoryId.value,
      amount: yuanToFen(amount.value),
      note: note.value.trim() || undefined,
      occurred_at: buildOccurredAt()
    })
    amount.value = null
    note.value = ''
    message.success('支出记录已保存')
    emit('created')
  } finally {
    loading.value = false
  }
}

onMounted(loadCategories)
</script>

