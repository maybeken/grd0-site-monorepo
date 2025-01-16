<template>
  <div v-if="loading" class="rounded-full overflow-hidden bg-gradient-to-r from-accent via-secondary to-accent blur-sm" :class="`h-${h} w-${w}`">
    <div class="bg-gradient-to-r from-secondary to-accent rounded-full h-full w-full motion-translate-x-loop-100 motion-duration-1000 motion-ease-spring-smooth"></div>
  </div>
  <slot v-else />
</template>

<script setup lang="ts">
const {
  loading = true,
  h = 'full',
  w = 'full',
} = defineProps<{
  loading: boolean,
  h: number | string,
  w: number | string,
}>();
</script>