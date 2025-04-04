<template>
  <img
    class="bg-secondary"
    :src="uri ? generateResizedSrc(uri, resolution, quality) : ''"
    :srcset="uri ? generateSrcset(uri, quality) : ''"
    loading="lazy"
  />
</template>

<script setup lang="ts">
import { generateResizedSrc, generateSrcset } from '@/helpers/cdn';

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