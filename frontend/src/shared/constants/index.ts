export const PERIOD_TYPE = { DAY: 1, WEEK: 2, MONTH: 3, YEAR: 4, LIFE: 5 } as const
export const SOURCE = { AI: 1, USER: 2 } as const
export const PERIOD_LABELS: Record<number, string> = { 1: '日报', 2: '周报', 3: '月报', 4: '年报', 5: '人生总结' }
export const PAGE_SIZE = 20
