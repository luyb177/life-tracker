import { ref, computed } from 'vue'
import { useMessage } from 'naive-ui'
import { useQuery } from '@tanstack/vue-query'
import { getExpenseByDate, createCategory, deleteCategory, ExpenseLogInfo } from '@/entities/expense/api/expense.api'
import { expenseKeys } from '@/entities/expense/api/expense.keys'
import { useExpenseCategories } from '@/entities/expense/model/useExpenseCategories'
import { useExpenseMutations } from '@/features/expense-create/useExpenseMutations'
import { formatDate } from '@/shared/utils/format'

export function useExpenses() {
  const msg = useMessage()
  const date = ref<number>(Date.now())
  const dateStr = computed(() => formatDate(new Date(date.value)))
  const showModal = ref(false)
  const editTarget = ref<ExpenseLogInfo | null>(null)
  const submitting = ref(false)
  const showCatModal = ref(false)
  const newCatName = ref('')

  const { data, isLoading } = useQuery({
    queryKey: computed(() => expenseKeys.byDate(dateStr.value)),
    queryFn: () => getExpenseByDate(dateStr.value),
  })
  const { categories } = useExpenseCategories()
  const { createMut, updateMut, deleteMut } = useExpenseMutations()

  const list = computed(() => ((data as any)?.value?.list || []) as ExpenseLogInfo[])
  const total = computed(() => (data as any)?.value?.total || 0)
  const expenseModalTitle = computed(() => editTarget.value ? '编辑支出' : '记一笔')

  function openCreate() { editTarget.value = null; showModal.value = true }
  function openEdit(e: ExpenseLogInfo) { editTarget.value = e; showModal.value = true }
  function closeEditor() { showModal.value = false }
  function openCategoryManager() { showCatModal.value = true }
  function closeCategoryManager() { showCatModal.value = false }

  async function handleSubmit(formData: { category_id: number; amount: number; note: string }) {
    submitting.value = true
    const occurred = dateStr.value + ' ' + new Date().toTimeString().slice(0, 8)
    try {
      if (editTarget.value) {
        await updateMut.mutateAsync({ id: editTarget.value.id, ...formData, occurred_at: occurred })
        msg.success('已更新')
      } else {
        await createMut.mutateAsync({ ...formData, occurred_at: occurred })
        msg.success('已记录')
      }
      closeEditor()
    } catch (e: any) { msg.error(e.response?.data?.msg || '操作失败') }
    finally { submitting.value = false }
  }

  async function handleDelete(id: number) {
    try { await deleteMut.mutateAsync(id); msg.success('已删除') }
    catch (e: any) { msg.error(e.response?.data?.msg || '删除失败') }
  }

  async function handleCreateCategory() {
    try { await createCategory(newCatName.value); newCatName.value = ''; msg.success('分类已创建') }
    catch (e: any) { msg.error(e.response?.data?.msg || '创建失败') }
  }

  async function handleDeleteCategory(id: number) {
    try { await deleteCategory(id); msg.success('分类已删除') }
    catch (e: any) { msg.error(e.response?.data?.msg || '删除失败') }
  }

  return {
    date,
    dateStr,
    showModal,
    editTarget,
    submitting,
    showCatModal,
    newCatName,
    isLoading,
    categories,
    list,
    total,
    expenseModalTitle,
    openCreate,
    openEdit,
    closeEditor,
    openCategoryManager,
    closeCategoryManager,
    handleSubmit,
    handleDelete,
    handleCreateCategory,
    handleDeleteCategory,
  }
}
