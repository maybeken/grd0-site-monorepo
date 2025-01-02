<script setup lang="ts">
import { listAssets, getGalleryCategory } from '@/services/gallery';
import { useGalleryStore } from '@/stores/gallery';

import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
dayjs.extend(relativeTime);

import type { Asset, AssetFileList, GalleryCategory } from '@/interfaces/Gallery';

const ASSET_URL = import.meta.env.VITE_ASSETS_DOMAIN;
const cdn_config = {
  resolution: 768,
  quality: 75,
};

const $store = useGalleryStore();
const files = listAssets();
const gallery_category = getGalleryCategory();

function getGalleryList(files: AssetFileList): AssetFileList {
  if (!files) return {};

  const list = Object.fromEntries(
    Object.entries(files).filter(
        ([key, val]) => key.indexOf('/gallery') === 0
    )
  );

  return list;
}

function getFileList(files: AssetFileList, filter: string): Asset[] {
  if (!files || !filter) return [];

  let assets: Asset[] = [];

  if (filter === 'all') {
    for (const category in files) {
      const assets_with_category = files[category].map((val) => { return {...val, category}; });

      assets = [...assets, ...assets_with_category];
    }
  } else {
    const filtered = files[filter] || [];
    assets = filtered.map((val) => { return {...val, category: filter}; })
  }

  assets = assets.sort((a, b) => { return dayjs(a.exif?.datetime).unix() - dayjs(b.exif?.datetime).unix() });

  return filter === 'all' ? assets.reverse() : assets;
}

function getCategoryList(list: AssetFileList | undefined): string[] {
  const default_category = 'all';

  if (!list) return [default_category];

  return [default_category, ...Object.keys(list)];
}

function formatCategoryName(value: string): string {
  const category_id = value.replace('/gallery/', '');

  if (gallery_category?.value && gallery_category?.value[category_id] && gallery_category?.value[category_id]?.title) {
    return gallery_category?.value[category_id].title;
  }

  return category_id.replace(/(^\w|\s\w)/g, (m: string) => m.toUpperCase())
}

function getCategoryCover(value: string): string | undefined {
  const category_id = value.replace('/gallery/', '');

  if (gallery_category?.value && gallery_category?.value[category_id]) {
    return gallery_category?.value[category_id].cover || undefined;
  }
}
</script>

<template>
  <div class="py-2">
    <DropdownSelection
      :options="getCategoryList(getGalleryList(files))"
      :selected="$store.selected_category"
      :stylize="formatCategoryName"
      @select="(newVal: string) => { $store.selected_category = newVal }"
    >
    </DropdownSelection>
  </div>
  <div class="py-2 relative rounded-xl" v-if="getCategoryCover($store.selected_category)">
    <img
      class="rounded-xl w-full object-cover object-center max-h-32 md:max-h-48 lg:max-h-96 blur-[2px] brightness-50"
      :src="`//${ASSET_URL}/cdn-cgi/image/width=${cdn_config.resolution},quality=${cdn_config.quality}${$store.selected_category}/${getCategoryCover($store.selected_category)}`"
    />
    <div class="absolute top-1/2 text-center w-full">
      <p class="px-4 md:text-2xl font-bold tracking-widest uppercase font-serif">{{ formatCategoryName($store.selected_category) }}</p>
    </div>
  </div>
  <div class="py-2 grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-2">
    <GalleryCard
      v-for="img of getFileList(getGalleryList(files), $store.selected_category)"
      :category="$store.selected_category"
      :image="img"
      :key="img.filename"
    ></GalleryCard>
  </div>
</template>
