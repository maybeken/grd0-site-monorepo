<template>
  <div
    class="border-solid border-accent border-4 rounded-full overflow-clip w-16 h-16" 
  >
    <img
      :src="$props.src"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

interface Props {
  src: string;
  size: number; // Does not work
};

const $props = defineProps<Props>();

// const defaultSize = 16;
</script>