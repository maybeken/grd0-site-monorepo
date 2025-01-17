<template>
  <div v-if="loading" class="rounded-full overflow-hidden bg-gradient-to-r from-accent via-secondary to-accent blur-sm" :class="[h_options[h], w_options[w]]">
    <div class="bg-gradient-to-r from-secondary to-accent rounded-full h-full w-full motion-translate-x-loop-100 motion-duration-1000 motion-ease-spring-smooth"></div>
  </div>
  <slot v-else />
</template>

<script setup lang="ts">
const h_options: {[key: string]: string} = {
  xs: 'h-4',
  sm: 'h-8',
  md: 'h-16',
  lg: 'h-24',
  xl: 'h-32',
};

const w_options: {[key: string]: string} = {
  xs: 'w-16',
  sm: 'w-24',
  md: 'w-48',
  lg: 'w-72',
  xl: 'w-96',
  '1/3': 'w-1/3',
  '2/3': 'w-2/3',
  full: 'w-full',
};

const {
  loading = true,
  h = 'md',
  w = 'full',
} = defineProps<{
  loading: boolean,
  h: string,
  w: string,
}>();
</script>