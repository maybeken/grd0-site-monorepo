<template>
  <GalleryZoom v-if="zoom_uri" :uri="zoom_uri" @close="zoom_uri = ''"></GalleryZoom>
  <template v-if="loading && asset_list.length === 0">
    <GalleryCard v-for="i in 10" :loading="true"></GalleryCard>
  </template>
  <template v-else>
    <GalleryCard
      v-for="img of asset_list"
      :image="img"
      :key="img.filename"
      :loading="false"
      @zoom="
        (uri: string) => {
          zoom_uri = uri
        }
      "
    ></GalleryCard>
  </template>
  <div v-if="loading && asset_list.length > 0" class="col-span-full flex justify-center py-4">
    <GalleryCard v-for="i in 5" :loading="true"></GalleryCard>
  </div>
  <div v-if="has_more" ref="sentinel" class="col-span-full h-1"></div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { loadMoreAssets } from '@/services/gallery'

import type { Asset } from '@/interfaces/Gallery'

const { collection = 'all' } = defineProps<{
  collection: string
}>()

const asset_list = ref<Asset[]>([])
const loading = ref(false)
const current_page = ref(1)
const has_more = ref(false)
const zoom_uri = ref('')
const sentinel = ref<HTMLElement | null>(null)

let observer: IntersectionObserver | null = null

async function fetchAssets(page: number) {
  if (loading.value) return
  loading.value = true
  try {
    const response = await loadMoreAssets(collection, page).send()
    if (page === 1) {
      asset_list.value = response.data
    } else {
      asset_list.value = [...asset_list.value, ...response.data]
    }
    current_page.value = page
    has_more.value = page < response.total_pages
  } finally {
    loading.value = false
  }
}

function setupObserver() {
  if (observer) observer.disconnect()
  const margin = window.innerWidth < 768 ? '400px' : '200px'
  observer = new IntersectionObserver(
    (entries) => {
      const entry = entries[0]
      if (entry?.isIntersecting && has_more.value && !loading.value) {
        fetchAssets(current_page.value + 1)
      }
    },
    { rootMargin: margin }
  )
  if (sentinel.value) observer.observe(sentinel.value)
}

watch(sentinel, () => setupObserver())

watch(() => collection, () => {
  asset_list.value = []
  current_page.value = 1
  has_more.value = false
  fetchAssets(1)
})

onMounted(() => {
  fetchAssets(1)
})

onUnmounted(() => observer?.disconnect())
</script>
