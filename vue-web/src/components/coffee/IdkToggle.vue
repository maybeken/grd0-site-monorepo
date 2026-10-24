<template>
  <label class="inline-flex cursor-pointer select-none items-center gap-1.5 text-sm font-medium">
    <input
      type="checkbox"
      :checked="modelValue"
      class="h-4 w-4 accent-accent"
      @change="$emit('update:modelValue', ($event.target as HTMLInputElement).checked)"
    />
    <span class="text-foreground/60">{{ $t('coffee.tastingForm.idk') }}</span>
  </label>
</template>

<script setup lang="ts">
interface Props {
  modelValue?: boolean
}

defineProps<Props>()
defineEmits<{
  'update:modelValue': [value: boolean]
}>()
</script>
