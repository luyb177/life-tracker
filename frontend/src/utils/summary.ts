import type { SummaryInfo } from '@/types/api'
import { addDays, formatDate, monthEndExclusive, monthStart, nextYearStart } from '@/utils/date'

export const summaryPeriodOptions = [
  { label: '日', value: 1 },
  { label: '周', value: 2 },
  { label: '月', value: 3 },
  { label: '年', value: 4 }
]

export function periodLabel(type: number) {
  return ['未知', '日报', '周报', '月报', '年报'][type] || '总结'
}

export function summaryPeriodTitle(summary: SummaryInfo) {
  if (summary.period_type === 1) return summary.period_start
  if (summary.period_type === 3) return summary.period_start.slice(0, 7)
  if (summary.period_type === 4) return summary.period_start.slice(0, 4)
  return summaryPeriodRange(summary)
}

export function summaryPeriodRange(summary: SummaryInfo) {
  return `${summary.period_start} 至 ${summary.period_end}`
}

export function currentPeriodStart(periodType: number) {
  const today = new Date()
  if (periodType === 3) return monthStart(today)
  if (periodType === 4) return `${today.getFullYear()}-01-01`
  if (periodType === 2) {
    const day = today.getDay() || 7
    today.setDate(today.getDate() - day + 1)
  }
  return formatDate(today)
}

export function normalizePeriodStart(periodType: number, dateText: string) {
  const date = new Date(`${dateText}T00:00:00`)
  if (periodType === 2) {
    const day = date.getDay() || 7
    date.setDate(date.getDate() - day + 1)
  }
  if (periodType === 3) {
    date.setDate(1)
  }
  if (periodType === 4) {
    date.setMonth(0, 1)
  }
  return formatDate(date)
}

export function periodEnd(periodType: number, start: string) {
  const startDate = new Date(`${start}T00:00:00`)
  if (periodType === 1) return addDays(start, 1)
  if (periodType === 2) return addDays(start, 7)
  if (periodType === 3) return monthEndExclusive(startDate)
  return nextYearStart(startDate)
}
