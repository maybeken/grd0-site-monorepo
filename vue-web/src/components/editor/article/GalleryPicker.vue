<template>
  <div class="flex flex-col gap-3 rounded-xl border border-dotted p-3">
    <div class="flex flex-wrap items-center gap-2">
      <span class="text-lg font-bold">{{ $t('article.editor.galleryTitle') }}</span>
      <select
        v-if="!$props.collection"
        v-model="selected"
        class="rounded border border-foreground/40 bg-background px-2.5 py-1.5 text-lg font-medium"
      >
        <option v-for="option in collections" :key="option.id" :value="option.id">
          {{ option.label }}
        </option>
      </select>
      <EditorButton variant="secondary" size="sm" class="ml-auto" @click="$emit('close')">
        {{ $t('common.cancel') }}
      </EditorButton>
    </div>

    <p v-if="!loading && assets.length === 0" class="text-lg font-medium text-foreground/60">
      {{ $t('article.editor.noAssets') }}
    </p>

    <div class="grid max-h-[50vh] grid-cols-2 gap-2 overflow-y-auto md:grid-cols-4">
      <button
        v-for="asset in assets"
        :key="assetUri(asset)"
        class="flex flex-col gap-1 rounded-xl bg-shade p-1 text-left text-lg font-medium hover:bg-accent/10 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-foreground/60"
        @click="select(asset)"
      >
        <CDNImage
          class="aspect-square w-full rounded-lg object-cover"
          :uri="assetUri(asset)"
          :resolution="320"
          :quality="70"
        ></CDNImage>
        <span class="truncate" :title="asset.filename">{{ asset.filename }}</span>
      </button>
    </div>

    <div class="flex justify-center">
      <EditorButton
        v-if="page < total_pages"
        size="sm"
        variant="secondary"
        :disabled="loading"
        @click="load(page + 1)"
      >
        {{ $t('common.loadMore') }}
      </EditorButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import EditorButton from '@/components/editor/ui/EditorButton.vue'
import CDNImage from '@/components/shared/CDNImage.vue'
import { getGalleryCollectionRaw, loadMoreAssets } from '@/services/gallery'

import type { Asset, GalleryCollection } from '@/interfaces/Gallery'

const $props = defineProps<{
  collection?: string
}>()

const $emit = defineEmits<{
  select: [asset: { collection: string; filename: string }]
  close: []
}>()

const { t } = useI18n()

const collections = ref<{ id: string; label: string }[]>([])
const selected = ref($props.collection ?? 'all')
const assets = ref<Asset[]>([])
const page = ref(0)
const total_pages = ref(0)
const loading = ref(false)

function assetUri(asset: Asset): string {
  if (asset.collection) return `${asset.collection}/${asset.filename}`

  return `/gallery/${selected.value}/${asset.filename}`
}

async function loadCollections(): Promise<void> {
  if ($props.collection) {
    collections.value = [{ id: $props.collection, label: $props.collection }]
    return
  }

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

watch(
  () => $props.collection,
  (value) => {
    selected.value = value ?? 'all'
  }
)

onMounted(() => {
  loadCollections()
  load(1)
})
</script>
