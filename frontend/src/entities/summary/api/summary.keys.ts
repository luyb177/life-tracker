export const summaryKeys = {
  all: ['summaries'] as const,
  day: (date: string) => [...summaryKeys.all, 'day', date] as const,
  range: (params: { start: string; end: string; period_type?: number }) => [...summaryKeys.all, 'range', params] as const,
  list: (params: { page_size: number; period_type?: number; page_token?: string }) => [...summaryKeys.all, 'list', params] as const,
  tagStats: (params: { start: string; end: string }) => [...summaryKeys.all, 'tagStats', params] as const,
  tagTrend: (params: { start: string; end: string }) => [...summaryKeys.all, 'tagTrend', params] as const,
}
