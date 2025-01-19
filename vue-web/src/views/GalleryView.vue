<script setup lang="ts">
import { ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getGalleryCategory } from '@/services/gallery';

import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
dayjs.extend(relativeTime);

const $route = useRoute();
const $router = useRouter();

const ASSET_URL = import.meta.env.VITE_ASSETS_URL;
const cdn_config = {
  resolution: 768,
  quality: 75,
};

const response = getGalleryCategory();
const gallery_category = response.data;
const loading = response.loading;

const uri = $route.params.category ?? 'all';
const selected_category = ref(typeof uri === "string" ? uri : uri[0]);

watch($route, (to, from) => {
  const uri = to.params.category ?? 'all';
  const uri_stripped = typeof uri === "string" ? uri : uri[0];
  selected_category.value = uri_stripped;
})

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

function selectCategory(category: string): void {
  if (category === 'all') {
    $router.push({ path: '/gallery' });
    return;
  }

  $router.push({ path: `${category}` });
}
</script>

<template>
  <DropdownSelection :disabled="loading" :options="gallery_category" :selected="selected_category"
    :stylize="formatCategoryName" @select="selectCategory">
  </DropdownSelection>
  <div class="py-2 relative rounded-xl" v-if="getCategoryCover(selected_category)">
    <img class="rounded-xl w-full object-cover object-center max-h-32 md:max-h-48 lg:max-h-96 blur-[2px] brightness-50"
      :src="`${ASSET_URL}/cdn-cgi/image/width=${cdn_config.resolution},quality=${cdn_config.quality}/gallery/${selected_category}/${getCategoryCover(selected_category)}`" />
    <div class="absolute top-1/2 text-center w-full">
      <p class="px-4 md:text-2xl font-bold tracking-widest uppercase font-serif">{{
        formatCategoryName(selected_category) }}</p>
    </div>
  </div>
  <div class="py-2 grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-2">
    <GalleryGrid :category="selected_category" :key="selected_category"></GalleryGrid>
  </div>
</template>
