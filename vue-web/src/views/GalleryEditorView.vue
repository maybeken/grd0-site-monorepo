<script setup lang="ts">
import { ref, watch } from 'vue';
import { listAssets } from '@/services/gallery';

import dayjs from 'dayjs';

import type { Ref } from 'vue';
import type { Asset, AssetFileList, GalleryDetail } from '@/interfaces/Gallery';

const ASSET_URL = import.meta.env.VITE_ASSETS_URL;
const cdn_config = {
  resolution: 768,
  quality: 75,
};

const files = listAssets();

const gallery_details: Ref<GalleryDetail[]> = ref([]);

function getGalleryList(files: AssetFileList): AssetFileList {
  if (!files) return {};

  const list = Object.fromEntries(
    Object.entries(files).filter(
        ([key, val]) => key.indexOf('/gallery') === 0
    )
  );

  return list;
}

function getFileList(files: AssetFileList): Asset[] {
  if (!files) return [];

  let assets: Asset[] = [];

  for (const category in files) {
    const assets_with_category = files[category].map((val) => { return {...val, category}; });

    assets = [...assets, ...assets_with_category];
  }

  return assets.reverse();
}
</script>

<template>
  <div class="py-2 grid grid-cols-1 grid-cols-3 gap-2">
    <div
      v-for="img of getFileList(getGalleryList(files))"
    >
      <p>{{ dayjs(img.exif?.datetime).format('YYYY/MM/DD HH:mm:ss') }}</p>
      <p class="capitalize">{{ img.category?.replace('/gallery/', '') }}</p>
      <img
        class="rounded-xl w-full object-cover aspect-square"
        :src="`${ASSET_URL}/cdn-cgi/image/width=${cdn_config.resolution},quality=${cdn_config.quality}${img.category}/${img.filename}`"
        loading="lazy"
      />
    </div>
  </div>
</template>
