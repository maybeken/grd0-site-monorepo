<template>
  <select
    :id="id"
    :value="modelValue"
    :name="name"
    :required="required"
    :aria-required="required ? 'true' : undefined"
    :aria-invalid="invalid ? 'true' : undefined"
    :disabled="disabled"
    class="w-full rounded border bg-background px-2.5 py-1.5 text-lg font-medium transition-colors focus:border-foreground/80 focus:outline focus:outline-2 focus:outline-offset-2 focus:outline-foreground/60 disabled:cursor-not-allowed disabled:opacity-50"
    :class="invalid ? 'border-red-400/70' : 'border-foreground/40'"
    @change="$emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
  >
    <slot></slot>
  </select>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    modelValue?: string | number | null | undefined
    id?: string
    name?: string
    required?: boolean
    invalid?: boolean
    disabled?: boolean
  }>(),
  {
    modelValue: '',
    id: '',
    name: '',
    required: false,
    invalid: false,
    disabled: false
  }
)

defineEmits<{
  'update:modelValue': [value: string]
}>()
</script>
