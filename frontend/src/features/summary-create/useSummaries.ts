import { ref, computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { useMessage } from 'naive-ui'
import { getSummaryList, generateAISummary, PERIOD_LABELS, SummaryInfo } from '@/entities/summary/api/summary.api'
import { summaryKeys } from '@/entities/summary/api/summary.keys'
import { useSummaryMutations } from '@/features/summary-create/useSummaryMutations'
import { todayStr, isAISummary, isUserRecord } from '@/shared/utils/format'
import { PAGE_SIZE } from '@/shared/constants'

export function useSummaries() {
  const msg = useMessage()
  const periodType = ref(1)
  const showEditor = ref(false)
  const editTarget = ref<SummaryInfo | null>(null)
  const submitting = ref(false)

  const { data, isLoading } = useQuery({
    queryKey: computed(() => summaryKeys.list({ page_size: PAGE_SIZE, period_type: periodType.value })),
    queryFn: () => getSummaryList({ page_size: PAGE_SIZE, period_type: periodType.value }),
  })

  const { createMut, updateMut, deleteMut } = useSummaryMutations()
  const list = computed(() => ((data as any)?.value?.list || []) as SummaryInfo[])

  const aiSummaries = computed(() => list.value.filter(s => isAISummary(s.source)))
  const userSummaries = computed(() => list.value.filter(s => isUserRecord(s.source)))
  const isEditing = computed(() => !!editTarget.value)
  const editorTitle = computed(() => isEditing.value ? '编辑总结' : `创建${PERIOD_LABELS[periodType.value]}`)

  function openCreate() {
    editTarget.value = null
    showEditor.value = true
  }

  function closeEditor() {
    showEditor.value = false
  }

  async function handleCreate(formData: any) {
    submitting.value = true
    try { await createMut.mutateAsync(formData); msg.success('已创建'); closeEditor() }
    catch (e: any) { msg.error(e.response?.data?.msg || '创建失败') }
    finally { submitting.value = false }
  }

  function openEdit(id: number) {
    const item = list.value.find(s => s.id === id)
    if (item) { editTarget.value = item; showEditor.value = true }
  }

  async function handleUpdate(formData: any) {
    submitting.value = true
    try { await updateMut.mutateAsync({ id: editTarget.value!.id, summary_content: formData.summary_content, tags: formData.tags, title: formData.title }); msg.success('已更新'); closeEditor() }
    catch (e: any) { msg.error(e.response?.data?.msg || '更新失败') }
    finally { submitting.value = false }
  }

  async function handleDelete(id: number) {
    try { await deleteMut.mutateAsync(id); msg.success('已删除') }
    catch (e: any) { msg.error(e.response?.data?.msg || '删除失败') }
  }

  async function handleGenerate() {
    try { await generateAISummary(periodType.value, todayStr()); msg.success(`AI ${PERIOD_LABELS[periodType.value]} 已生成`) }
    catch (e: any) { msg.error(e.response?.data?.msg || '生成失败') }
  }

  return {
    periodType,
    showEditor,
    editTarget,
    submitting,
    isEditing,
    editorTitle,
    isLoading,
    list,
    aiSummaries,
    userSummaries,
    openCreate,
    openEdit,
    closeEditor,
    handleCreate,
    handleUpdate,
    handleDelete,
    handleGenerate,
  }
}
