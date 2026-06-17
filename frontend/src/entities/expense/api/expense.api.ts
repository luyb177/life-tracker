import client from '@/shared/api/client'

export interface ExpenseCategoryInfo {
  id: number
  name: string
  type: number
}

export interface ExpenseCategoryListData {
  categories: ExpenseCategoryInfo[]
}

export interface ExpenseLogInfo {
  id: number; category: { id: number; name: string; type: number }
  amount: number; note: string; location: string; occurred_at: string; created_at: string
}
export function getExpenseByDate(date: string) { return client.get('/expense/by_date', { params: { date } }).then(r => r.data.data) }
export function getExpenseList(params: { page_size: number; page_token?: string }) { return client.get('/expense/list', { params }).then(r => r.data.data) }
export function createExpense(data: { category_id: number; amount: number; note?: string; occurred_at: string }) { return client.post('/expense/create', data).then(r => r.data.data) }
export function updateExpense(data: { id: number; category_id?: number; amount?: number; note?: string; occurred_at?: string }) { return client.post('/expense/update', data).then(r => r.data.data) }
export function deleteExpense(id: number) { return client.post('/expense/delete', { id }).then(r => r.data.data) }
export function getCategories() { return client.get('/expense/categories').then(r => r.data.data as ExpenseCategoryListData) }
export function createCategory(name: string) { return client.post('/expense/category/create', { name }).then(r => r.data.data) }
export function deleteCategory(id: number) { return client.post('/expense/category/delete', { id }).then(r => r.data.data) }
export function getDailyTotal(date: string) { return client.get('/expense/daily_total', { params: { date } }).then(r => r.data.data) }
export function getStatsRange(s: string, e: string) { return client.get('/expense/stats/range', { params: { start: s, end: e } }).then(r => r.data.data) }
export function getStatsCategory(s: string, e: string) { return client.get('/expense/stats/category', { params: { start: s, end: e } }).then(r => r.data.data) }
export function getStatsTrend(s: string, e: string) { return client.get('/expense/stats/trend', { params: { start: s, end: e } }).then(r => r.data.data) }
export function getMonthlyTrend(s: string, e: string) { return client.get('/expense/stats/monthly', { params: { start: s, end: e } }).then(r => r.data.data) }
