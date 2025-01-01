<script setup lang="ts">
import { listAssets, getGalleryCategory } from '@/services/gallery';
import { useGalleryStore } from '@/stores/gallery';

import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
dayjs.extend(relativeTime);

import type { Asset, AssetFileList, GalleryCategory } from '@/interfaces/Gallery';

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

  if (gallery_category?.value && gallery_category?.value[category_id]) {
    return gallery_category?.value[category_id];
  }

  return category_id.replace(/(^\w|\s\w)/g, (m: string) => m.toUpperCase())
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
  <div class="py-2 grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-2">
    <GalleryCard
      v-for="img of getFileList(getGalleryList(files), $store.selected_category)"
      :category="$store.selected_category"
      :image="img"
      :key="img.filename"
    ></GalleryCard>
  </div>
</template>
