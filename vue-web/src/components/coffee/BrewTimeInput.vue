<template>
  <input
    type="text"
    :value="isFocused ? rawInput : displayValue"
    @input="onInput"
    @focus="onFocus"
    @blur="onBlur"
    :disabled="disabled"
    placeholder="mm:ss"
    class="bg-background border border-accent rounded px-2 py-1 w-20 text-center disabled:opacity-50"
  />
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

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

const isFocused = ref(false)
const rawInput = ref('')

const displayValue = computed(() => {
  if (props.modelValue == null) return ''
  const mins = Math.floor(props.modelValue / 60)
  const secs = props.modelValue % 60
  return `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`
})

function parseValue(input: string): number | null {
  const cleaned = input.replace(/[^0-9:]/g, '')
  if (cleaned === '') return null

  const parts = cleaned.split(':')
  if (parts.length === 2) {
    const mins = parseInt(parts[0] || '0') || 0
    const secs = parseInt(parts[1] || '0') || 0
    return mins * 60 + secs
  } else if (parts.length === 1 && parts[0] !== '') {
    const digits = (parts[0] ?? '').padStart(4, '0').slice(-4)
    const mins = parseInt(digits.slice(0, 2)) || 0
    const secs = parseInt(digits.slice(2, 4)) || 0
    return mins * 60 + secs
  }
  return null
}

function onFocus() {
  isFocused.value = true
  rawInput.value = displayValue.value
}

function onInput(e: Event) {
  const input = (e.target as HTMLInputElement).value.replace(/[^0-9:]/g, '')
  rawInput.value = input

  if (input.includes(':')) {
    emit('update:modelValue', parseValue(input))
  }
}

function onBlur() {
  isFocused.value = false
  emit('update:modelValue', parseValue(rawInput.value))
}
</script>
