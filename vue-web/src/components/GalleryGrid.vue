<template>
  <template v-if="loading">
    <GalleryCard v-for="i in 10" :loading="true"></GalleryCard>
  </template>
  <template v-else>
    <GalleryCard v-for="img of asset_list" :image="img" :key="img.filename" :loading="loading"></GalleryCard>
  </template>
</template>

<script setup lang="ts">
import { listAssets } from '@/services/gallery';

const { category = 'all' } = defineProps<{
  category: string,
}>();

let response = listAssets(category);
const asset_list = response?.data;
const loading = response?.loading;
</script>