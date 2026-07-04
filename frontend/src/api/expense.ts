import { http, unwrap } from '@/api/client'
import type {
  CategoryStat,
  CursorList,
  ExpenseCategoryInfo,
  ExpenseLogInfo,
  IDResponse,
  MonthTrendPoint,
  TrendPoint
} from '@/types/api'

export interface CreateExpensePayload {
  category_id: number
  amount: number
  note?: string
  occurred_at: string
}

export function createExpense(payload: CreateExpensePayload) {
  return unwrap<IDResponse>(http.post('/expense/create', payload))
}

export function listExpenseCategories() {
  return unwrap<{ categories: ExpenseCategoryInfo[] }>(http.get('/expense/categories'))
}

export function createExpenseCategory(name: string) {
  return unwrap<Record<string, never>>(http.post('/expense/category/create', { name }))
}

export function refundExpense(id: number) {
  return unwrap<Record<string, never>>(http.post('/expense/refund', { id }))
}

export function getExpensesByDate(date: string) {
  return unwrap<{ list: ExpenseLogInfo[]; total: number }>(
    http.get('/expense/by_date', { params: { date } })
  )
}

export function listExpenses(params: { page_size: number; page_token?: string }) {
  return unwrap<CursorList<ExpenseLogInfo>>(http.get('/expense/list', { params }))
}

export function getCategoryStats(params: { start: string; end: string }) {
  return unwrap<{ categories: CategoryStat[] }>(http.get('/expense/stats/category', { params }))
}

export function getTrendStats(params: { start: string; end: string }) {
  return unwrap<{ points: TrendPoint[] }>(http.get('/expense/stats/trend', { params }))
}

export function getMonthlyTrendStats(params: { start: string; end: string }) {
  return unwrap<{ points: MonthTrendPoint[] }>(http.get('/expense/stats/monthly', { params }))
}
