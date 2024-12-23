<script setup lang="ts">
import { ref, watch } from 'vue';
import { listAssets } from '@/services/gallery';
import { useGalleryStore } from '@/stores/gallery';

const $store = useGalleryStore();
const files = listAssets();

function getGalleryList(files, filter) {
  if (!files) return;

  const list = Object.fromEntries(
    Object.entries(files).filter(
        ([key, val]) => key.indexOf('/gallery') === 0
    )
  );

  if (filter) {
    return list[filter];
  }

  return list ?? [];
}

function getCategoryList(list) {
  if (!list) return;

  return Object.keys(list);
}

// Set default category
watch(files, (newVal) => {
  if (!$store.selected_category.value) {
    $store.selected_category = getCategoryList(getGalleryList(newVal))[0];
  }
})
</script>

<template>
  <div class="py-2">
    <DropdownSelection
      :options="getCategoryList(getGalleryList(files))"
      :selected="$store.selected_category"
      :stylize="value => value.replace('/gallery/', '').replace(/(^\w|\s\w)/g, m => m.toUpperCase())"
      @select="(newVal) => { $store.selected_category = newVal }"
    >
    </DropdownSelection>
  </div>
  <div class="py-2 grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-2">
    <GalleryCard
      :selected_category="$store.selected_category"
      :list="getGalleryList(files, $store.selected_category)"
    ></GalleryCard>
  </div>
</template>
