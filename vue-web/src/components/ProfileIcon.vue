<template>
  <div
    class="border-solid border-accent border-4 rounded-full overflow-clip w-16" 
  >
    <img
      :src="$props.src"
      alt="Profile Picture"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

interface Props {
  src: string;
};

const $props = defineProps<Props>();
</script>