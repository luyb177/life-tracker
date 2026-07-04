import { BarChart3, BookOpenText, CalendarDays, ReceiptText, Settings } from '@lucide/vue'

export const navItems = [
  { label: '今天', path: '/today', icon: CalendarDays },
  { label: '记录', path: '/life-logs', icon: BookOpenText },
  { label: '支出', path: '/expenses', icon: ReceiptText },
  { label: '总结', path: '/summaries', icon: BarChart3 },
  { label: '设置', path: '/settings', icon: Settings }
]

