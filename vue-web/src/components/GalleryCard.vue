<template>
  <div class="flex flex-col gap-2 justify-end p-1 pb-2 rounded-xl bg-shade text-secondary text-sm">
    <div>
      <CDNImage
        class="rounded-xl w-full object-cover aspect-square lg:cursor-zoom-in"
        :resolution="768"
        :quality="75"
        :uri="$props.image?.filename ? `${$props.image?.category}/${$props.image?.filename}` : ''"
        @click="zoomImage(`${$props.image?.category}/${$props.image?.filename}`)"
      ></CDNImage>
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

import { getGalleryDetail } from '@/services/gallery';

import type { Asset, GalleryDetail } from '@/interfaces/Gallery';
import Skeleton from './Skeleton.vue';

interface Props {
  image?: Asset,
  loading: boolean,
}

const $props = defineProps<Props>();

const details = ref<GalleryDetail>();
const response = getGalleryDetail($props.image?.category);
const gallery_details = response?.data || ref([]);

const $emit = defineEmits<{
  zoom: [uri?: string]
}>();

watch(gallery_details, (newVal) => {
  if (!newVal) return;
  
  const file_rule = newVal.find((item) => item.filename === $props.image?.filename);
  const category_rule = newVal.find((item) => item.filename === '*');

  if (category_rule) {
    details.value = category_rule;

    if (file_rule) {
      details.value = file_rule;
    }
  } else if (file_rule) {
    details.value = file_rule;
  }
});

function zoomImage(val: string) {
  $emit('zoom', val);
};
</script>