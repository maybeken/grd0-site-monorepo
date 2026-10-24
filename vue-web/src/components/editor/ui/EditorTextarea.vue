<template>
  <textarea
    :id="id"
    :value="modelValue"
    :name="name"
    :rows="rows"
    :required="required"
    :aria-required="required ? 'true' : undefined"
    :aria-invalid="invalid ? 'true' : undefined"
    :placeholder="placeholder"
    :disabled="disabled"
    class="w-full resize-y rounded border bg-background px-2.5 py-1.5 text-lg font-medium leading-relaxed placeholder:text-foreground/60 transition-colors focus:border-foreground/80 focus:outline focus:outline-2 focus:outline-offset-2 focus:outline-foreground/60 disabled:cursor-not-allowed disabled:opacity-50"
    :class="invalid ? 'border-red-400/70' : 'border-foreground/40'"
    @input="$emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
  ></textarea>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    modelValue?: string | number | null | undefined
    id?: string
    name?: string
    rows?: number
    required?: boolean
    invalid?: boolean
    placeholder?: string
    disabled?: boolean
  }>(),
  {
    modelValue: '',
    id: '',
    name: '',
    rows: 3,
    required: false,
    invalid: false,
    placeholder: '',
    disabled: false
  }
)

defineEmits<{
  'update:modelValue': [value: string]
}>()
</script>
