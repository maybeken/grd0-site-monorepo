<template>
  <div class="rounded-xl bg-shade" v-for="img of $props.list">
    <a :href="`//${ASSET_URL}${selected_category}/${img.filename}`" target="_blank">
      <img class="p-1 rounded-xl w-full aspect-square"
        :src="`//${ASSET_URL}/cdn-cgi/image/width=512,quality=75${selected_category}/${img.filename}`" lazy />
    </a>
    <div class="p-2 pt-1 text-secondary text-xs">
      <p class="text-center">{{ img.exif?.equipment?.camera }} {{ img.exif?.equipment?.lens }}</p>
      <p class="text-center">ISO {{ img.exif?.iso }} | {{ img.exif?.fstop }} | {{ img.exif?.shutter }}s</p>
      <p class="pt-2 text-right">{{ img.exif?.datetime }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Asset } from '@/interfaces/Gallery';

const ASSET_URL = import.meta.env.VITE_ASSETS_DOMAIN;

interface Props {
  selected_category: string,
  list?: Asset[];
}

const $props = defineProps<Props>();
</script>