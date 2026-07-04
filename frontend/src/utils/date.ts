export function formatDate(date = new Date()) {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

export function startOfDay(date = new Date()) {
  return new Date(date.getFullYear(), date.getMonth(), date.getDate())
}

export function daysBetween(targetText: string, base = new Date()) {
  const target = startOfDay(new Date(`${targetText}T00:00:00`)).getTime()
  const today = startOfDay(base).getTime()
  return Math.round((target - today) / 86400000)
}

export function relativeDateLabel(dateText: string) {
  const diff = daysBetween(dateText)
  if (diff === 0) return '今天'
  if (diff === -1) return '昨天'
  if (diff === -2) return '前天'
  return dateText
}

export function isFutureDateTimestamp(timestamp: number) {
  return timestamp > startOfDay().getTime()
}

export function formatDateTime(date = new Date()) {
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${formatDate(date)} ${hours}:${minutes}:${seconds}`
}

export function addDays(dateText: string, days: number) {
  const date = new Date(`${dateText}T00:00:00`)
  date.setDate(date.getDate() + days)
  return formatDate(date)
}

export function monthStart(date = new Date()) {
  return formatDate(new Date(date.getFullYear(), date.getMonth(), 1))
}

export function monthEndExclusive(date = new Date()) {
  return formatDate(new Date(date.getFullYear(), date.getMonth() + 1, 1))
}

export function yearStart(date = new Date()) {
  return formatDate(new Date(date.getFullYear(), 0, 1))
}

export function nextYearStart(date = new Date()) {
  return formatDate(new Date(date.getFullYear() + 1, 0, 1))
}
