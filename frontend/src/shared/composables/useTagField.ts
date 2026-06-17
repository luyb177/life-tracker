import { ref, computed } from 'vue'
import { parseTags, joinTags } from '@/shared/utils/format'

export function useTagField(initial: string = '') {
  const raw = ref(initial)
  const list = computed({ get: () => parseTags(raw.value), set: (v: string[]) => { raw.value = joinTags(v) } })
  return { raw, list }
}
