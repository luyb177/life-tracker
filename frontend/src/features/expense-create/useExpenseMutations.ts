import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { createExpense, updateExpense, deleteExpense } from '@/entities/expense/api/expense.api'
import { expenseKeys } from '@/entities/expense/api/expense.keys'

export function useExpenseMutations() {
  const qc = useQueryClient()
  const invalidate = () => qc.invalidateQueries({ queryKey: expenseKeys.all })

  const createMut = useMutation({ mutationFn: createExpense, onSuccess: invalidate })
  const updateMut = useMutation({ mutationFn: updateExpense, onSuccess: invalidate })
  const deleteMut = useMutation({ mutationFn: (id: number) => deleteExpense(id), onSuccess: invalidate })

  return { createMut, updateMut, deleteMut }
}
