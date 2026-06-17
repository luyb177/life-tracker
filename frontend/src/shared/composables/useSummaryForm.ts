import { computed, ref, watch } from 'vue'
import { PERIOD_TYPE, PERIOD_LABELS } from '@/shared/constants'
import { alignPeriodStart, calcPeriodEnd, formatDate } from '@/shared/utils/format'
import { useTagField } from '@/shared/composables/useTagField'

export interface SummaryFormOptions {
  initialContent?: string
  initialTags?: string
  initialTitle?: string
  initialPeriodType?: number
  initialPeriodStart?: string
  initialPeriodEnd?: string
  lockPeriodFields?: boolean
}

export interface SummaryFormPayload {
  period_type: number
  period_start: string
  period_end: string
  summary_content: string
  tags: string
  title?: string
}

export function useSummaryForm(options: SummaryFormOptions = {}) {
  const periodType = ref(options.initialPeriodType || PERIOD_TYPE.DAY)
  const ts = ref(options.initialPeriodStart ? new Date(options.initialPeriodStart).getTime() : Date.now())
  const periodStart = computed(() => formatDate(new Date(ts.value)))
  const periodEnd = computed(() => calcPeriodEnd(periodType.value, periodStart.value))
  const content = ref(options.initialContent || '')
  const title = ref(options.initialTitle || '')
  const { raw: tags } = useTagField(options.initialTags || '')
  const isLife = computed(() => periodType.value === PERIOD_TYPE.LIFE)
  const lifeEndTs = ref(options.initialPeriodEnd ? new Date(options.initialPeriodEnd).getTime() : Date.now() + 86400000)
  const lifeEndStr = computed(() => formatDate(new Date(lifeEndTs.value)))
  const lifeError = ref(false)
  const periodLocked = computed(() => !!options.lockPeriodFields)
  const displayPeriodEnd = computed(() => {
    if (periodLocked.value && options.initialPeriodEnd) return options.initialPeriodEnd
    return isLife.value ? lifeEndStr.value : periodEnd.value
  })
  const displayPeriodLabel = computed(() => PERIOD_LABELS[periodType.value] || '总结')
  const ptOptions = Object.entries(PERIOD_LABELS).map(([k, v]) => ({ label: v, value: Number(k) }))

  watch(periodType, (newType) => {
    if (!periodLocked.value && newType !== PERIOD_TYPE.DAY) {
      ts.value = alignPeriodStart(newType, new Date(ts.value)).getTime()
    }
  })

  watch(ts, (newTs) => {
    if (!periodLocked.value && periodType.value !== PERIOD_TYPE.DAY && periodType.value !== PERIOD_TYPE.LIFE) {
      const aligned = alignPeriodStart(periodType.value, new Date(newTs))
      if (aligned.getTime() !== newTs) ts.value = aligned.getTime()
    }
    lifeError.value = isLife.value && lifeEndTs.value <= newTs
  })

  watch(lifeEndTs, (endTs) => {
    lifeError.value = isLife.value && endTs <= ts.value
  })

  function buildPayload(): SummaryFormPayload | null {
    if (isLife.value && lifeEndTs.value <= ts.value) {
      lifeError.value = true
      return null
    }

    return {
      period_type: periodType.value,
      period_start: periodStart.value,
      period_end: isLife.value ? lifeEndStr.value : periodEnd.value,
      summary_content: content.value,
      tags: tags.value,
      title: title.value || undefined,
    }
  }

  return {
    periodType,
    ts,
    periodStart,
    periodEnd,
    content,
    title,
    tags,
    isLife,
    lifeEndTs,
    lifeEndStr,
    lifeError,
    periodLocked,
    displayPeriodEnd,
    displayPeriodLabel,
    ptOptions,
    buildPayload,
  }
}
