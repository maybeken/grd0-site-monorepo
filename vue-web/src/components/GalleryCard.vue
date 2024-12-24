<template>
  <div class="rounded-xl bg-shade">
    <a :href="`//${ASSET_URL}${$props.category}/${image.filename}`" target="_blank">
      <img class="p-1 rounded-xl w-full aspect-square"
        :src="`//${ASSET_URL}/cdn-cgi/image/width=512,quality=75${$props.category}/${$props.image.filename}`" lazy />
    </a>
    <div class="p-2 pt-1 text-secondary text-xs">
      <p class="pb-2 text-center">{{ details?.description }}</p>
      <p class="text-center">{{ $props.image.exif?.equipment?.camera }} {{ $props.image.exif?.equipment?.lens }}</p>
      <p class="text-center">ISO {{ $props.image.exif?.iso }} | {{ $props.image.exif?.fstop }} | {{ $props.image.exif?.shutter }}s</p>
      <p class="pt-2 text-right">{{ details?.tz_adjustment ? dayjs($props.image.exif?.datetime).add(details?.tz_adjustment, 'h').format('LLL') : dayjs($props.image.exif?.datetime).format('LLL') }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import dayjs from 'dayjs';
import LocalizedFormat from 'dayjs/plugin/localizedFormat';
dayjs.extend(LocalizedFormat);

import { getGalleryDetail } from '@/services/gallery';

import type { Asset, GalleryDetail } from '@/interfaces/Gallery';

const gallery_details = getGalleryDetail();

const ASSET_URL = import.meta.env.VITE_ASSETS_DOMAIN;

interface Props {
  category: string,
  image: Asset;
}

const $props = defineProps<Props>();
const details = ref();

function findGalleryDetail(gallery_details: GalleryDetail[], path: string, filename: string): GalleryDetail | null {
  if (!gallery_details) return null;

  return gallery_details.find((val) => {
    return val.path === path && val.filename === filename;
  }) ?? null;
}

watch(gallery_details, (newVal: GalleryDetail[]) => {
  if (!newVal) return;

  details.value = findGalleryDetail(newVal, $props.category, $props.image.filename);
})
</script>