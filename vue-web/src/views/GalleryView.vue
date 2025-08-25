<script setup lang="ts">
import { ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getGalleryCollection } from '@/services/gallery';
import { formatCollectionName, getCollectionCover } from '@/helpers/collection';

import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
dayjs.extend(relativeTime);

const $route = useRoute();
const $router = useRouter();

const response = getGalleryCollection();
const gallery_collection = response.data;
const loading = response.loading;

const uri = $route.params.collection ?? 'all';

let selected_collection = ref('');

if (typeof uri === 'string') {
  selected_collection.value = uri;
} else if (Array.isArray(uri) && uri.length > 0) {
  selected_collection.value = uri[0]!;
}

watch($route, (to, from) => {
  const uri = to.params.collection ?? 'all';

  let selected_collection = ref('');

  if (typeof uri === 'string') {
    selected_collection.value = uri;
  } else if (Array.isArray(uri) && uri.length > 0) {
    selected_collection.value = uri[0]!;
  }
})

function selectCollection(collection: string): void {
  if (collection === 'all') {
    $router.push({ path: '/gallery' });
    return;
  }

  $router.push({ path: `${collection}` });
}
</script>

<template>
  <div class="pb-4">
    <p class="text-2xl font-bold">Captures The Moment<CursorBlink /></p>
  </div>
  <DropdownSelection :disabled="loading" :options="gallery_collection" :selected="selected_collection"
    :stylize="(val: string) => formatCollectionName(gallery_collection, val)" @select="selectCollection">
  </DropdownSelection>
  <div class="py-2 relative rounded-xl" v-if="getCollectionCover(gallery_collection, selected_collection)">
    <CDNImage
      class="rounded-xl w-full object-cover object-center max-h-32 md:max-h-48 lg:max-h-96 blur-[2px] brightness-50"
      :resolution="768"
      :quality="75"
      :uri="`/gallery/${selected_collection}/${getCollectionCover(gallery_collection, selected_collection)}`"
    ></CDNImage>
    <div class="absolute top-1/2 text-center w-full">
      <p class="px-4 md:text-2xl font-bold tracking-widest uppercase font-serif">{{
        formatCollectionName(gallery_collection, selected_collection) }}</p>
    </div>
  </div>
  <div class="py-2 grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-2">
    <GalleryGrid :collection="selected_collection" :key="selected_collection"></GalleryGrid>
  </div>
</template>
