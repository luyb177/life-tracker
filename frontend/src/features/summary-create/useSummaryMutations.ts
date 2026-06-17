import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { createSummary, updateSummary, deleteSummary } from '@/entities/summary/api/summary.api'
import { summaryKeys } from '@/entities/summary/api/summary.keys'

export function useSummaryMutations() {
  const qc = useQueryClient()
  const invalidate = () => qc.invalidateQueries({ queryKey: summaryKeys.all })

  const createMut = useMutation({ mutationFn: createSummary, onSuccess: invalidate })
  const updateMut = useMutation({ mutationFn: updateSummary, onSuccess: invalidate })
  const deleteMut = useMutation({ mutationFn: (id: number) => deleteSummary(id), onSuccess: invalidate })

  return { createMut, updateMut, deleteMut }
}
