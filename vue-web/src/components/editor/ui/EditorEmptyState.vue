<template>
  <div class="flex flex-col items-start gap-2 py-8">
    <p class="text-lg font-medium text-foreground/60">{{ message }}</p>
    <p v-if="description" class="text-sm font-medium leading-relaxed text-foreground/60">
      {{ description }}
    </p>
    <slot></slot>
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
