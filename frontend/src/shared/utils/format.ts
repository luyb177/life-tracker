import { PERIOD_LABELS } from '@/entities/summary/api/summary.api'

export function formatDate(d: string | Date): string {
  const date = typeof d === 'string' ? new Date(d) : d
  return date.toISOString().slice(0, 10)
}

export function todayStr(): string { return formatDate(new Date()) }

export function yesterdayStr(): string {
  const d = new Date(); d.setDate(d.getDate() - 1); return formatDate(d)
}

export function calcPeriodEnd(periodType: number, start: string): string {
  const d = new Date(start)
  switch (periodType) {
    case 1: d.setDate(d.getDate() + 1); break
    case 2: d.setDate(d.getDate() + 7); break
    case 3: d.setMonth(d.getMonth() + 1); break
    case 4: d.setFullYear(d.getFullYear() + 1); break
    default: d.setDate(d.getDate() + 1)
  }
  return formatDate(d)
}

export function periodLabel(t: number): string { return PERIOD_LABELS[t] || '总结' }

export function parseTags(tags: string): string[] { return tags ? tags.split(',').map(t => t.trim()).filter(Boolean) : [] }
export function joinTags(tags: string[]): string { return tags.join(',') }
export function isAISummary(source: number): boolean { return source === 1 }
export function isUserRecord(source: number): boolean { return source === 2 }

export function formatPeriodRange(start: string, end: string): string {
  return `${start} ~ ${end}`
}

export function alignToMonday(d: Date): Date {
  const day = d.getDay()
  const diff = day === 0 ? -6 : 1 - day
  const monday = new Date(d); monday.setDate(d.getDate() + diff)
  return monday
}

export function alignToMonthStart(d: Date): Date {
  return new Date(d.getFullYear(), d.getMonth(), 1)
}

export function alignToYearStart(d: Date): Date {
  return new Date(d.getFullYear(), 0, 1)
}

export function alignPeriodStart(periodType: number, d: Date): Date {
  switch (periodType) {
    case 2: return alignToMonday(d)
    case 3: return alignToMonthStart(d)
    case 4: return alignToYearStart(d)
    default: return d
  }
}
