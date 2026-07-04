export function formatDate(date = new Date()) {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
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

