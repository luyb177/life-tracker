import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { expenseKeys } from '@/entities/expense/api/expense.keys'
import { getCategories } from '@/entities/expense/api/expense.api'

export function useExpenseCategories() {
  const query = useQuery({
    queryKey: expenseKeys.categories(),
    queryFn: () => getCategories(),
  })

  const categories = computed(() => query.data.value?.categories || [])
  const options = computed(() => categories.value.map(category => ({
    label: category.name,
    value: category.id,
  })))

  return {
    ...query,
    categories,
    options,
  }
}
