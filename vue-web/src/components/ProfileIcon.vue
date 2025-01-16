<template>
  <div
    class="border-solid border-accent border-4 rounded-full overflow-clip w-16 bg-secondary"
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