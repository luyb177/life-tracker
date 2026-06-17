<script setup lang="ts">
import { NButton, NInput, NTag } from 'naive-ui'
import ConfirmAction from '@/shared/ui/confirm-action/ConfirmAction.vue'

defineProps<{
  categories: Array<{ id: number; name: string; type: number }>
  newCategoryName: string
}>()
const emit = defineEmits<{
  'update:newCategoryName': [value: string]
  create: []
  delete: [id: number]
}>()
</script>
<template>
  <div style="padding: 16px">
    <div style="display: flex; gap: 8px; margin-bottom: 16px">
      <NInput :value="newCategoryName" placeholder="新分类" style="flex: 1" @update:value="emit('update:newCategoryName', $event)" />
      <NButton size="small" @click="emit('create')">添加</NButton>
    </div>
    <div v-for="c in categories" :key="c.id" style="display: flex; justify-content: space-between; align-items: center; padding: 8px 0; border-bottom: 1px solid var(--n-border-color)">
      <span>{{ c.name }} <NTag v-if="c.type === 1" size="tiny">系统</NTag></span>
      <ConfirmAction v-if="c.type === 2" @confirm="emit('delete', c.id)" />
    </div>
  </div>
</template>
