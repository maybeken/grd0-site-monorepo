<template>
  <img
    class="bg-secondary"
    :src="`${ASSET_URL}/cdn-cgi/image/width=${resolution},quality=${quality}${uri}`"
    loading="lazy"
  />
</template>

<script setup lang="ts">
interface Props {
  uri: string,
  resolution: number;
  quality: number;
};

const ASSET_URL = import.meta.env.VITE_ASSETS_URL;
const {
  uri,
  resolution = 768,
  quality = 75,
} = defineProps<Props>();
</script>