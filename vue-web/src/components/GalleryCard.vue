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
    <div v-if="gallery_details?.description || loading" class="px-1">
      <Skeleton h="md" w="full" :loading="loading">
        <p class="text-justify">{{ gallery_details?.description }}</p>
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
        <p class="text-right">{{ gallery_details?.tz_adjustment ?
          dayjs($props.image?.exif?.datetime).add(gallery_details?.tz_adjustment, 'h').format('LLL') :
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

import { getGalleryDetail } from '@/services/gallery';

import type { Asset } from '@/interfaces/Gallery';
import Skeleton from './Skeleton.vue';

const ASSET_URL = import.meta.env.VITE_ASSETS_URL;

interface Props {
  image?: Asset,
  loading: boolean,
}

const $props = defineProps<Props>();
const details = ref();

const response = getGalleryDetail($props.image?.category);
const gallery_details = response?.data;

const cdn_config = {
  resolution: 768,
  quality: 75,
};
</script>