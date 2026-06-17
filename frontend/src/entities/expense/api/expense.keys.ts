export const expenseKeys = {
  all: ['expenses'] as const,
  byDate: (date: string) => [...expenseKeys.all, 'byDate', date] as const,
  list: (params: { page_size: number; page_token?: string }) => [...expenseKeys.all, 'list', params] as const,
  categories: () => [...expenseKeys.all, 'categories'] as const,
  dailyTotal: (date: string) => [...expenseKeys.all, 'dailyTotal', date] as const,
  statsRange: (s: string, e: string) => [...expenseKeys.all, 'statsRange', s, e] as const,
  statsCategory: (s: string, e: string) => [...expenseKeys.all, 'statsCategory', s, e] as const,
  statsTrend: (s: string, e: string) => [...expenseKeys.all, 'statsTrend', s, e] as const,
  monthlyTrend: (s: string, e: string) => [...expenseKeys.all, 'monthlyTrend', s, e] as const,
}
