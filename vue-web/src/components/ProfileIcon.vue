<template>
  <div
    class="border-solid border-accent border-4 rounded-full overflow-clip h-full w-full bg-secondary"
  >
    <img
      v-if="src"
      class="motion-preset-fade"
      :src="src"
    />
  </div>
</template>

<script setup lang="ts">
interface Props {
  src: string;
};

const {
  src,
} = defineProps<Props>();
</script>