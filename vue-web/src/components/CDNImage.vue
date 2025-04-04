<template>
  <img
    class="bg-secondary"
    :src="uri ? `${ASSET_URL}/cdn-cgi/image/fit=scale-down,width=${resolution},quality=${quality},format=${format}${uri}` : ''"
    :srcset="uri ? generateSrcset(uri, quality) : ''"
    loading="lazy"
  />
</template>

<script setup lang="ts">
import { generateSrcset } from '@/helpers/cdn';

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
const format = 'jpeg';
</script>