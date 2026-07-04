import { http, unwrap } from '@/api/client'
import type { CursorList, IDResponse, LifeLogInfo, TagInfo } from '@/types/api'

export interface CreateLifeLogPayload {
  content: string
  occurred_at: string
  tags?: TagInfo[]
}

export function createLifeLog(payload: CreateLifeLogPayload) {
  return unwrap<IDResponse>(http.post('/life_log/create', payload))
}

export function getLifeLogsByDate(date: string) {
  return unwrap<{ list: LifeLogInfo[] }>(http.get('/life_log/by_date', { params: { date } }))
}

export function listLifeLogs(params: { page_size: number; page_token?: string }) {
  return unwrap<CursorList<LifeLogInfo>>(http.get('/life_log/list', { params }))
}

export function deleteLifeLog(id: number) {
  return unwrap<Record<string, never>>(http.post('/life_log/delete', { id }))
}

