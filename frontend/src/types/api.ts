export interface ApiEnvelope<T> {
  code: number
  msg: string
  data: T
}

export interface UserInfo {
  user_id: number
  username: string
  email: string
  avatar: string
  created_at: string
}

export interface TagInfo {
  id: number
  name: string
}

export interface LoginResp {
  token: string
  refresh_token: string
  user_info: UserInfo
}

export interface IDResponse {
  id: number
}

export interface LifeLogInfo {
  id: number
  content: string
  tags?: TagInfo[]
  occurred_at: string
  created_at: string
}

export interface ExpenseCategoryInfo {
  id: number
  name: string
  type: number
}

export interface ExpenseLogInfo {
  id: number
  category: ExpenseCategoryInfo
  amount: number
  note: string
  location?: string
  occurred_at: string
  status: number
  refunded_at?: string
  created_at: string
}

export interface SummaryInfo {
  id: number
  period_type: number
  period_start: string
  period_end: string
  source: number
  summary_content: string
  suggestion_content?: string
  title?: string
  tags?: TagInfo[]
  location?: string
  created_at: string
  updated_at: string
}

export interface CursorList<T> {
  list: T[]
  page_token: string
  has_more: boolean
}

export interface CategoryStat {
  category_id: number
  category_name: string
  total: number
}

export interface TrendPoint {
  date: string
  total: number
}

export interface MonthTrendPoint {
  month: string
  total: number
}
