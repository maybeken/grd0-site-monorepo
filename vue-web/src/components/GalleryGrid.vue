<template>
  <GalleryZoom v-if="zoom_uri" :uri="zoom_uri" @close="zoom_uri = ''"></GalleryZoom>
  <template v-if="loading">
    <GalleryCard v-for="i in 10" :loading="true"></GalleryCard>
  </template>
  <template v-else>
    <GalleryCard v-for="img of asset_list" :image="img" :key="img.filename" :loading="loading" @zoom="(uri: string) => { zoom_uri = uri; }"></GalleryCard>
  </template>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { listAssets } from '@/services/gallery';

const { category = 'all' } = defineProps<{
  category: string,
}>();

let response = listAssets(category);
const asset_list = response?.data;
const loading = response?.loading;

const zoom_uri = ref('');
</script>