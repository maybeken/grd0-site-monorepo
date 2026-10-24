<template>
  <div class="flex flex-col gap-1.5">
    <div class="flex flex-wrap items-center justify-between gap-2">
      <label v-if="label" :for="htmlFor" class="text-lg font-medium leading-snug">
        {{ label }}<span v-if="required" class="text-red-400" aria-hidden="true"> *</span>
      </label>
      <slot name="label-action"></slot>
    </div>

    <slot></slot>

    <p v-if="error" class="text-sm font-medium text-red-400">{{ error }}</p>
    <p v-else-if="help" class="text-sm font-medium leading-relaxed text-foreground/60">{{ help }}</p>
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    label?: string
    htmlFor?: string
    required?: boolean
    error?: string
    help?: string
  }>(),
  {
    label: '',
    htmlFor: '',
    required: false,
    error: '',
    help: ''
  }
)
</script>
