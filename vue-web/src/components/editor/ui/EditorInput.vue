<template>
  <input
    ref="input"
    :id="id"
    :type="type"
    :value="modelValue"
    :name="name"
    :required="required"
    :aria-required="required ? 'true' : undefined"
    :aria-invalid="invalid ? 'true' : undefined"
    :placeholder="placeholder"
    :disabled="disabled"
    :autocomplete="autocomplete"
    class="w-full rounded border bg-background px-2.5 py-1.5 text-lg font-medium placeholder:text-foreground/60 transition-colors focus:border-foreground/80 focus:outline focus:outline-2 focus:outline-offset-2 focus:outline-foreground/60 disabled:cursor-not-allowed disabled:opacity-50"
    :class="invalid ? 'border-red-400/70' : 'border-foreground/40'"
    @input="onInput"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue?: string | number | null | undefined
    id?: string
    type?: string
    name?: string
    required?: boolean
    invalid?: boolean
    placeholder?: string
    disabled?: boolean
    autocomplete?: string
    number?: boolean
  }>(),
  {
    modelValue: '',
    id: '',
    type: 'text',
    name: '',
    required: false,
    invalid: false,
    placeholder: '',
    disabled: false,
    autocomplete: 'off',
    number: false
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: string | number | null]
}>()

const input = ref<HTMLInputElement | null>(null)

function onInput(event: Event): void {
  const raw = (event.target as HTMLInputElement).value

  if (props.number) {
    emit('update:modelValue', raw === '' ? null : Number(raw))
    return
  }

  emit('update:modelValue', raw)
}

defineExpose({
  focus: () => input.value?.focus()
})
</script>
