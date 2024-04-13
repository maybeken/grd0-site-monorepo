<template>
  <div
    class="border-solid border-accent border-4 rounded-full" 
    :class="$props.size ? [`h-${$props.size}`, `w-${$props.size}`] : [`h-${defaultSize}`,`w-${defaultSize}`]"
  >
    <img
      class="rounded-full"
      :class="$props.size ? [`h-${$props.size}`, `w-${$props.size}`] : [`h-${defaultSize}`,`w-${defaultSize}`]"
      :src="$props.src"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

interface Props {
  src: string;
  size: number;
};

const $props = defineProps<Props>();

const defaultSize = 16;
</script>