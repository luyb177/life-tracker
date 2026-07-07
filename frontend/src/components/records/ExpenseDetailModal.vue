<template>
  <n-modal :show="show" preset="dialog" title="支出详情" class="detail-modal" @update:show="emit('update:show', $event)">
    <n-form class="detail-form" :show-label="true" label-placement="top">
      <n-form-item label="发生时间">
        <n-date-picker
          v-model:formatted-value="form.occurred_at"
          type="datetime"
          value-format="yyyy-MM-dd HH:mm:ss"
          class="full-control"
        />
      </n-form-item>
      <div class="detail-two-columns">
        <n-form-item label="金额">
          <n-input-number v-model:value="form.amount" :min="0" :precision="2" class="full-control">
            <template #prefix>¥</template>
          </n-input-number>
        </n-form-item>
        <n-form-item label="分类">
          <n-select
            v-model:value="form.category_id"
            :options="categoryOptions"
            :loading="loadingCategories"
            class="full-control"
          />
        </n-form-item>
      </div>
      <n-form-item label="备注">
        <n-input v-model:value="form.note" />
      </n-form-item>
      <p v-if="expense" class="detail-meta">
        创建 {{ expense.created_at }} · 更新 {{ expense.last_updated_at || expense.updated_at }}
      </p>
      <p v-if="expense?.location" class="detail-meta">地点 {{ expense.location }}</p>
    </n-form>
    <template #action>
      <n-button tertiary type="error" :loading="loading" :disabled="!expense" @click="remove">删除</n-button>
      <n-button v-if="expense?.status === 0" tertiary type="warning" :loading="loading" @click="refund">
        退款
      </n-button>
      <n-button @click="emit('update:show', false)">取消</n-button>
      <n-button type="primary" :loading="loading" :disabled="!canSave" @click="save">保存</n-button>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { ExpenseCategoryInfo, ExpenseLogInfo } from '@/types/api'
import { fenToYuan, yuanToFen } from '@/utils/money'

const props = defineProps<{
  show: boolean
  expense?: ExpenseLogInfo | null
  categories: ExpenseCategoryInfo[]
  loading?: boolean
  loadingCategories?: boolean
}>()
const emit = defineEmits<{
  'update:show': [value: boolean]
  save: [payload: { id: number; category_id: number; amount: number; note: string; occurred_at: string }]
  delete: [id: number]
  refund: [id: number]
}>()

const form = ref({ amount: 0, category_id: null as number | null, note: '', occurred_at: '' })
const categoryOptions = computed(() => props.categories.map((item) => ({ label: item.name, value: item.id })))
const canSave = computed(() => Boolean(props.expense && form.value.amount > 0 && form.value.category_id))

watch(
  () => [props.expense, props.show] as const,
  ([expense, show]) => {
    if (!expense || !show) return
    form.value = {
      amount: fenToYuan(expense.amount),
      category_id: expense.category.id,
      note: expense.note,
      occurred_at: expense.occurred_at
    }
  },
  { immediate: true }
)

function save() {
  if (!props.expense || !form.value.category_id || !canSave.value) return
  emit('save', {
    id: props.expense.id,
    category_id: form.value.category_id,
    amount: yuanToFen(form.value.amount),
    note: form.value.note.trim(),
    occurred_at: form.value.occurred_at
  })
}

function remove() {
  if (!props.expense) return
  emit('delete', props.expense.id)
}

function refund() {
  if (!props.expense) return
  emit('refund', props.expense.id)
}
</script>
