<template>
  <input
    type="text"
    :value="displayValue"
    @input="onInput"
    :disabled="disabled"
    placeholder="mm:ss"
    class="bg-background border border-accent rounded px-2 py-1 w-20 text-center disabled:opacity-50"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  modelValue?: number | null
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false
})
const emit = defineEmits<{
  'update:modelValue': [value: number | null]
}>()

const displayValue = computed(() => {
  if (props.modelValue == null) return ''
  const mins = Math.floor(props.modelValue / 60)
  const secs = props.modelValue % 60
  return `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`
})

function onInput(e: Event) {
  const input = (e.target as HTMLInputElement).value.replace(/[^0-9:]/g, '')
  const parts = input.split(':')
  if (parts.length === 2) {
    const mins = parseInt(parts[0] || '0') || 0
    const secs = parseInt(parts[1] || '0') || 0
    emit('update:modelValue', mins * 60 + secs)
  } else if (parts.length === 1 && parts[0] !== '') {
    const secs = parseInt(parts[0] || '0') || 0
    emit('update:modelValue', secs)
  } else {
    emit('update:modelValue', null)
  }
}
</script>
