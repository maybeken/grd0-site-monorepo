<template>
  <div class="flex flex-col gap-2 justify-end p-1 pb-2 rounded-xl bg-shade text-secondary text-xs">
    <div>
      <a :href="`//${ASSET_URL}${$props.category}/${image.filename}`" target="_blank">
        <img class="rounded-xl w-full aspect-square cursor-zoom-in"
          :src="`//${ASSET_URL}/cdn-cgi/image/width=${cdn_config.resolution},quality=${cdn_config.quality}${$props.category}/${$props.image.filename}`" lazy />
      </a>
    </div>
    <div v-if="details?.description" class="px-1">
      <p class="text-justify">{{ details?.description }}</p>
    </div>
    <div class="px-1 grow">
      <p class="text-center">{{ $props.image.exif?.equipment?.camera }} {{ $props.image.exif?.equipment?.lens }}</p>
      <p class="text-center">ISO {{ $props.image.exif?.iso }} | {{ $props.image.exif?.fstop }} | {{ $props.image.exif?.shutter }}s</p>
    </div>
    <div class="px-1">
      <p class="text-right">{{ details?.tz_adjustment ? dayjs($props.image.exif?.datetime).add(details?.tz_adjustment, 'h').format('YYYY/MM/DD hh:mm A') : dayjs($props.image.exif?.datetime).format('LLL') }}</p>
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

const cdn_config = {
  resolution: 768,
  quality: 75,
};

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