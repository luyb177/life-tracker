<script setup lang="ts">
import { NSelect } from 'naive-ui'
import { useExpenseCategories } from '@/entities/expense/model/useExpenseCategories'

const props = defineProps<{ modelValue?: number | null }>()
const emit = defineEmits<{ 'update:modelValue': [v: number | null] }>()

const { options } = useExpenseCategories()
</script>
<template>
  <NSelect :value="props.modelValue" :options="options" placeholder="选择分类" @update:value="emit('update:modelValue', $event)" />
</template>
