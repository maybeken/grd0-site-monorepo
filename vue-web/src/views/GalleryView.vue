<script setup lang="ts">
import { ref } from 'vue';
import { getGalleryCategory } from '@/services/gallery';
import { useGalleryStore } from '@/stores/gallery';

import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
dayjs.extend(relativeTime);

const ASSET_URL = import.meta.env.VITE_ASSETS_URL;
const cdn_config = {
  resolution: 768,
  quality: 75,
};

const $store = useGalleryStore();
const response = getGalleryCategory();
const gallery_category = response.data;
const loading = response.loading;

const selected = ref($store.selected_category);

function formatCategoryName(value?: string): string | void {
  if (!value) return;

  const category_id = value.replace('/gallery/', '');
  let category_name = category_id;

  if (gallery_category?.value && gallery_category?.value[category_id]) {
    category_name = gallery_category?.value[category_id].title || category_name;
  }

  return category_name.replace(/(^\w|\s\w)/g, (m: string) => m.toUpperCase())
}

function getCategoryCover(value: string): string | undefined {
  const category_id = value.replace('/gallery/', '');

  if (gallery_category?.value && gallery_category?.value[category_id]) {
    return gallery_category?.value[category_id].cover || undefined;
  }
}
</script>

<template>
  <DropdownSelection :disabled="loading" :options="gallery_category" :selected="$store.selected_category"
    :stylize="formatCategoryName" @select="(newVal: string) => { $store.selected_category = newVal; selected = newVal; }">
  </DropdownSelection>
  <div class="py-2 relative rounded-xl" v-if="getCategoryCover($store.selected_category)">
    <img class="rounded-xl w-full object-cover object-center max-h-32 md:max-h-48 lg:max-h-96 blur-[2px] brightness-50"
      :src="`${ASSET_URL}/cdn-cgi/image/width=${cdn_config.resolution},quality=${cdn_config.quality}${$store.selected_category}/${getCategoryCover($store.selected_category)}`" />
    <div class="absolute top-1/2 text-center w-full">
      <p class="px-4 md:text-2xl font-bold tracking-widest uppercase font-serif">{{
        formatCategoryName($store.selected_category) }}</p>
    </div>
  </div>
  <div class="py-2 grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-2">
    <GalleryGrid :category="selected" :key="selected"></GalleryGrid>
  </div>
</template>
