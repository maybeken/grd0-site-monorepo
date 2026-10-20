<template>
  <div class="flex flex-col gap-3 rounded-xl border border-dotted p-3">
    <div class="flex flex-wrap items-center gap-2">
      <span class="font-bold">{{ $t('article.editor.galleryTitle') }}</span>
      <select v-model="selected" class="rounded-full border border-dotted px-3 py-1 bg-background">
        <option v-for="option in collections" :key="option.id" :value="option.id">
          {{ option.label }}
        </option>
      </select>
      <button class="ml-auto rounded-full border border-dotted px-3 py-1" @click="$emit('close')">
        {{ $t('common.cancel') }}
      </button>
    </div>

    <p v-if="!loading && assets.length === 0" class="text-sm text-secondary">
      {{ $t('article.editor.noAssets') }}
    </p>

    <div class="grid grid-cols-2 md:grid-cols-4 gap-2 max-h-[50vh] overflow-y-auto">
      <button
        v-for="asset in assets"
        :key="assetUri(asset)"
        class="flex flex-col gap-1 rounded-xl bg-shade p-1 text-left text-xs"
        @click="select(asset)"
      >
        <CDNImage
          class="w-full aspect-square object-cover rounded-lg"
          :uri="assetUri(asset)"
          :resolution="320"
          :quality="70"
        ></CDNImage>
        <span class="truncate">{{ asset.filename }}</span>
      </button>
    </div>

    <div class="flex justify-center">
      <button
        v-if="page < total_pages"
        class="rounded-full border border-dotted px-4 py-1 disabled:opacity-50"
        :disabled="loading"
        @click="load(page + 1)"
      >
        {{ $t('article.editor.loadMore') }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'

import { t } from '@/i18n'
import { getGalleryCollectionRaw, loadMoreAssets } from '@/services/gallery'

import type { Asset, GalleryCollection } from '@/interfaces/Gallery'

const $emit = defineEmits<{
  select: [asset: { collection: string; filename: string }]
  close: []
}>()

const collections = ref<{ id: string; label: string }[]>([])
const selected = ref('all')
const assets = ref<Asset[]>([])
const page = ref(0)
const total_pages = ref(0)
const loading = ref(false)

function assetUri(asset: Asset): string {
  if (asset.collection) return `${asset.collection}/${asset.filename}`

  return `/gallery/${selected.value}/${asset.filename}`
}

async function loadCollections(): Promise<void> {
  try {
    const response = await getGalleryCollectionRaw().send()
    const data = response as GalleryCollection
    if (!data || typeof data !== 'object' || Array.isArray(data)) return

    const options = Object.entries(data).map(([id, value]) => ({
      id,
      label: value?.title || id
    }))

    collections.value = [{ id: 'all', label: t('article.editor.allCollections') }, ...options]
  } catch {
    collections.value = [{ id: 'all', label: t('article.editor.allCollections') }]
  }
}

async function load(target: number): Promise<void> {
  loading.value = true

  try {
    const response = await loadMoreAssets(selected.value, target).send()
    const data = Array.isArray(response?.data) ? response.data : []

    assets.value = target === 1 ? data : [...assets.value, ...data]
    page.value = target
    total_pages.value = Number(response?.total_pages ?? 0)
  } catch {
    if (target === 1) assets.value = []
  } finally {
    loading.value = false
  }
}

function select(asset: Asset): void {
  let collection = asset.collection || ''

  if (!collection && selected.value !== 'all') {
    collection = `/gallery/${selected.value}`
  }
  if (!collection) return

  $emit('select', { collection, filename: asset.filename })
}

watch(selected, () => {
  assets.value = []
  page.value = 0
  total_pages.value = 0
  load(1)
})

onMounted(() => {
  loadCollections()
  load(1)
})
</script>
