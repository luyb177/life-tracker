import { computed, ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { useMessage } from 'naive-ui'
import { getDaySummaries, SummaryInfo } from '@/entities/summary/api/summary.api'
import { summaryKeys } from '@/entities/summary/api/summary.keys'
import { useSummaryMutations } from '@/features/summary-create/useSummaryMutations'
import { PERIOD_TYPE } from '@/shared/constants'
import { isAISummary, isUserRecord } from '@/shared/utils/format'

export function useRecords() {
  const msg = useMessage()
  const date = ref<number>(Date.now())
  const dateStr = computed(() => new Date(date.value).toISOString().slice(0, 10))
  const showEditor = ref(false)
  const editTarget = ref<SummaryInfo | null>(null)
  const submitting = ref(false)

  const { data, isLoading } = useQuery({
    queryKey: computed(() => summaryKeys.day(dateStr.value)),
    queryFn: () => getDaySummaries(dateStr.value),
  })

  const { createMut, updateMut, deleteMut } = useSummaryMutations()
  const list = computed(() => ((data.value as any)?.list || []) as SummaryInfo[])
  const userRecords = computed(() => list.value.filter(item => isUserRecord(item.source)))
  const aiSummaries = computed(() => list.value.filter(item => isAISummary(item.source)))
  const editorTitle = computed(() => editTarget.value ? '编辑记录' : '写记录')

  function openCreate() {
    editTarget.value = null
    showEditor.value = true
  }

  function openEdit(id: number) {
    const item = list.value.find(summary => summary.id === id)
    if (item) {
      editTarget.value = item
      showEditor.value = true
    }
  }

  function closeEditor() {
    showEditor.value = false
  }

  async function handleSubmit(formData: any) {
    submitting.value = true
    try {
      if (editTarget.value) {
        await updateMut.mutateAsync({
          id: editTarget.value.id,
          summary_content: formData.summary_content,
          tags: formData.tags,
          title: formData.title,
        })
        msg.success('已更新')
      } else {
        await createMut.mutateAsync({
          ...formData,
          period_type: PERIOD_TYPE.DAY,
          period_start: dateStr.value,
          period_end: dateStr.value,
        })
        msg.success('已创建')
      }
      closeEditor()
    } catch (e: any) {
      msg.error(e.response?.data?.msg || '操作失败')
    } finally {
      submitting.value = false
    }
  }

  async function handleDelete(id: number) {
    try {
      await deleteMut.mutateAsync(id)
      msg.success('已删除')
    } catch (e: any) {
      msg.error(e.response?.data?.msg || '删除失败')
    }
  }

  return {
    date,
    dateStr,
    showEditor,
    editTarget,
    submitting,
    isLoading,
    list,
    userRecords,
    aiSummaries,
    editorTitle,
    openCreate,
    openEdit,
    closeEditor,
    handleSubmit,
    handleDelete,
  }
}
