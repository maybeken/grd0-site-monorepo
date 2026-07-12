<template>
  <label class="inline-flex items-center gap-1 text-xs cursor-pointer select-none">
    <input
      type="checkbox"
      :checked="modelValue"
      @change="$emit('update:modelValue', ($event.target as HTMLInputElement).checked)"
      class="accent-accent w-3 h-3"
    />
    <span class="opacity-70">IDK</span>
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
