<script setup lang="ts">
import { ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getGalleryCategory } from '@/services/gallery';
import { formatCategoryName, getCategoryCover } from '@/helpers/category';

import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
dayjs.extend(relativeTime);

const $route = useRoute();
const $router = useRouter();

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

function selectCategory(category: string): void {
  if (category === 'all') {
    $router.push({ path: '/gallery' });
    return;
  }

  $router.push({ path: `${category}` });
}
</script>

<template>
  <div class="pb-4">
    <p class="text-2xl font-bold">Captures The Moment<CursorBlink /></p>
  </div>
  <DropdownSelection :disabled="loading" :options="gallery_category" :selected="selected_category"
    :stylize="(val: string) => formatCategoryName(gallery_category, val)" @select="selectCategory">
  </DropdownSelection>
  <div class="py-2 relative rounded-xl" v-if="getCategoryCover(gallery_category, selected_category)">
    <CDNImage
      class="rounded-xl w-full object-cover object-center max-h-32 md:max-h-48 lg:max-h-96 blur-[2px] brightness-50"
      :resolution="768"
      :quality="75"
      :uri="`/gallery/${selected_category}/${getCategoryCover(gallery_category, selected_category)}`"
    ></CDNImage>
    <div class="absolute top-1/2 text-center w-full">
      <p class="px-4 md:text-2xl font-bold tracking-widest uppercase font-serif">{{
        formatCategoryName(gallery_category, selected_category) }}</p>
    </div>
  </div>
  <div class="py-2 grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-2">
    <GalleryGrid :category="selected_category" :key="selected_category"></GalleryGrid>
  </div>
</template>
