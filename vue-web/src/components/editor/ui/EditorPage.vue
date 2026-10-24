<template>
  <div class="mx-auto flex w-full flex-col gap-6" :class="full ? '' : 'max-w-5xl'">
    <slot></slot>
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    full?: boolean
  }>(),
  {
    full: false
  }
)
</script>
