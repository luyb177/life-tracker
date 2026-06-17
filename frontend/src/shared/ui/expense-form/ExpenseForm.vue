<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { NInput, NButton, NSelect, NInputNumber } from 'naive-ui'
import { useExpenseCategories } from '@/entities/expense/model/useExpenseCategories'

const props = defineProps<{
  editId?: number | null
  initialCategoryId?: number | null
  initialAmount?: number | null
  initialNote?: string
  submitting?: boolean
}>()

const emit = defineEmits<{
  submit: [data: { category_id: number; amount: number; note: string }]
  cancel: []
}>()

const { options: catOptions } = useExpenseCategories()
const catId = ref<number | null>(props.initialCategoryId || null)
const amount = ref<number | null>(props.initialAmount || null)
const note = ref(props.initialNote || '')

onMounted(() => {
  if (props.editId && props.initialCategoryId) catId.value = props.initialCategoryId
  if (props.editId && props.initialAmount) amount.value = props.initialAmount
  if (props.editId && props.initialNote) note.value = props.initialNote
})

function handleSubmit() {
  if (!catId.value || !amount.value) return
  emit('submit', { category_id: catId.value, amount: amount.value, note: note.value })
}
</script>
<template>
  <div>
    <NSelect v-model:value="catId" :options="catOptions" placeholder="选择分类" style="margin-bottom: 12px" />
    <NInputNumber v-model:value="amount" placeholder="金额" :min="0" style="width: 100%; margin-bottom: 12px" />
    <NInput v-model:value="note" placeholder="备注（可选）" />
    <div style="display: flex; gap: 8px; margin-top: 12px; justify-content: flex-end">
      <NButton @click="emit('cancel')">取消</NButton>
      <NButton type="primary" :loading="submitting" @click="handleSubmit">{{ editId ? '保存' : '提交' }}</NButton>
    </div>
  </div>
</template>
