import client from '@/shared/api/client'

export interface SummaryInfo {
  id: number; period_type: number; period_start: string; period_end: string
  source: number; summary_content: string; suggestion_content: string
  title: string; tags: string; location: string; created_at: string; updated_at: string
}
export const PERIOD_LABELS: Record<number, string> = { 1: '日报', 2: '周报', 3: '月报', 4: '年报', 5: '人生总结' }
export const SOURCE_LABELS: Record<number, string> = { 1: 'AI', 2: '用户' }

export function getDaySummaries(date: string) { return client.get('/summary/day', { params: { date } }).then(r => r.data.data) }
export function getRangeSummaries(params: { start: string; end: string; period_type?: number }) { return client.get('/summary/range', { params }).then(r => r.data.data) }
export function getSummaryList(params: { page_size: number; period_type?: number; page_token?: string }) { return client.get('/summary/list', { params }).then(r => r.data.data) }
export function createSummary(data: { period_type: number; period_start: string; period_end: string; summary_content: string; title?: string; tags?: string }) { return client.post('/summary/create', data).then(r => r.data.data) }
export function updateSummary(data: { id: number; summary_content?: string; title?: string; tags?: string }) { return client.post('/summary/update', data).then(r => r.data.data) }
export function deleteSummary(id: number) { return client.post('/summary/delete', { id }).then(r => r.data.data) }
export function generateAISummary(period_type: number, period_start: string) { return client.post('/summary/generate', { period_type, period_start }).then(r => r.data.data) }
export function getTagStats(params: { start: string; end: string }) { return client.get('/summary/stats/tags', { params }).then(r => r.data.data) }
export function getTagTrend(params: { start: string; end: string }) { return client.get('/summary/stats/tag_trend', { params }).then(r => r.data.data) }
