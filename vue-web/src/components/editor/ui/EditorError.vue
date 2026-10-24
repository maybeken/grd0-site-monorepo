<template>
  <div
    class="flex flex-wrap items-center justify-between gap-4 rounded-xl border border-red-500/40 bg-red-500/5 p-4"
    role="alert"
  >
    <div class="flex min-w-0 flex-col gap-1">
      <p class="text-lg font-medium text-red-400">{{ message }}</p>
      <p v-if="description" class="text-sm font-medium leading-relaxed text-foreground/60">
        {{ description }}
      </p>
    </div>

    <div v-if="$slots.actions" class="flex flex-wrap items-center gap-2">
      <slot name="actions"></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    message: string
    description?: string
  }>(),
  {
    description: ''
  }
)
</script>
