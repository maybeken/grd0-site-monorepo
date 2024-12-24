<script setup lang="ts">
import { ref, watch } from 'vue';
import { listAssets } from '@/services/gallery';
import { useGalleryStore } from '@/stores/gallery';

import type { Asset, AssetFileList } from '@/interfaces/Gallery';

const $store = useGalleryStore();
const files = listAssets();

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
  if (!files) return [];

  return files[filter] || [];
}

function getCategoryList(list: AssetFileList | undefined): string[] {
  if (!list) return [];

  return Object.keys(list);
}

// Set default category
watch(files, (newVal) => {
  if (newVal && !$store.selected_category) {
    const full_list = getCategoryList(getGalleryList(newVal));

    if (full_list) $store.selected_category = full_list[0];
  }
})
</script>

<template>
  <div class="py-2">
    <DropdownSelection
      :options="getCategoryList(getGalleryList(files))"
      :selected="$store.selected_category"
      :stylize="(value: string) => value.replace('/gallery/', '').replace(/(^\w|\s\w)/g, (m: string) => m.toUpperCase())"
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
