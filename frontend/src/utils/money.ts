export function yuanToFen(value: number | null) {
  if (!value) return 0
  return Math.round(value * 100)
}

export function formatYuan(fen = 0) {
  return `¥${(fen / 100).toFixed(2)}`
}

