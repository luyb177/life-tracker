import { http, unwrap } from '@/api/client'
import type { CursorList, IDResponse, SummaryInfo, TagInfo } from '@/types/api'

export interface CreateSummaryPayload {
  period_type: number
  period_start: string
  period_end: string
  summary_content: string
  suggestion_content?: string
  title?: string
  tags?: TagInfo[]
}

export function listSummaries(params: {
  period_type?: number
  page_size: number
  page_token?: string
}) {
  return unwrap<CursorList<SummaryInfo>>(http.get('/summary/list', { params }))
}

export function rangeSummaries(params: { period_type?: number; start: string; end: string }) {
  return unwrap<{ list: SummaryInfo[] }>(http.get('/summary/range', { params }))
}

export function createSummary(payload: CreateSummaryPayload) {
  return unwrap<IDResponse>(http.post('/summary/create', payload))
}

export function generateAISummary(period_type: number, period_start: string) {
  return unwrap<SummaryInfo>(http.post('/summary/generate', { period_type, period_start }))
}

