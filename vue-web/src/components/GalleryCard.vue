<template>
  <div class="flex flex-col gap-2 justify-end p-1 pb-2 rounded-xl bg-shade text-secondary text-xs">
    <div>
      <a :href="$props.image?.filename ? `${ASSET_URL}${$props.image?.category}/${$props.image?.filename}` : ''"
        target="_blank">
        <img class="rounded-xl w-full object-cover aspect-square cursor-zoom-in bg-secondary"
          :src="$props.image?.filename ? `${ASSET_URL}/cdn-cgi/image/width=${cdn_config.resolution},quality=${cdn_config.quality}${$props.image?.category}/${$props.image?.filename}` : ''"
          loading="lazy" />
      </a>
    </div>
    <div v-if="details?.description || loading" class="px-1">
      <Skeleton h="md" w="full" :loading="loading">
        <p class="text-justify">{{ details?.description }}</p>
      </Skeleton>
    </div>
    <div class="px-1 grow">
      <Skeleton h="sm" w="full" :loading="loading">
        <p class="text-center">{{ $props.image?.exif?.equipment?.camera }} {{ $props.image?.exif?.equipment?.lens }}</p>
      </Skeleton>
      <Skeleton h="sm" w="full" :loading="loading">
        <p class="text-center"></p>
        <p class="text-center">ISO {{ $props.image?.exif?.iso }} | {{ $props.image?.exif?.fstop }} | {{
          $props.image?.exif?.shutter }}s</p>
      </Skeleton>
    </div>
    <div class="px-1">
      <Skeleton class="ml-auto" h="sm" w="2/3" :loading="loading">
        <p class="text-center"></p>
        <p class="text-right">{{ details?.tz_adjustment ?
          dayjs($props.image?.exif?.datetime).add(details?.tz_adjustment, 'h').format('LLL') :
          dayjs($props.image?.exif?.datetime).format('LLL') }}</p>
      </Skeleton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import dayjs from 'dayjs';
import LocalizedFormat from 'dayjs/plugin/localizedFormat';
dayjs.extend(LocalizedFormat);

import { listGallery } from '@/services/gallery';

import type { Asset, GalleryDetail } from '@/interfaces/Gallery';
import Skeleton from './Skeleton.vue';

const gallery_details = listGallery();

const ASSET_URL = import.meta.env.VITE_ASSETS_URL;

interface Props {
  image?: Asset,
  loading: boolean,
}

const $props = defineProps<Props>();
const details = ref();

const cdn_config = {
  resolution: 768,
  quality: 75,
};

function findGalleryDetail(gallery_details: GalleryDetail[], path: string, filename: string): GalleryDetail | null {
  if (!gallery_details) return null;

  const file_detail = gallery_details.find((val) => {
    return val.path === path && val.filename === filename;
  });

  const folder_detail = gallery_details.find((val) => {
    return val.path === path && val.filename === "*";
  });

  if (file_detail && folder_detail) {
    return { ...folder_detail, ...file_detail };
  } else if (file_detail) {
    return file_detail;
  } else if (folder_detail) {
    return folder_detail;
  }

  return null;
}

watch(gallery_details, (newVal: GalleryDetail[]) => {
  if (!newVal || !$props.image?.category) return;

  details.value = findGalleryDetail(newVal, $props.image?.category, $props.image.filename);
})
</script>