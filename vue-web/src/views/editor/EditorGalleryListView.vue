<template>
  <div class="flex flex-col gap-6">
    <header class="flex flex-col gap-2">
      <h1 class="text-2xl font-bold">{{ $t('editor.gallery') }}<CursorBlink /></h1>
      <p class="text-secondary">{{ $t('gallery.editor.listHint') }}</p>
    </header>

    <p v-if="loading" class="py-4 text-sm text-foreground/40">{{ $t('common.loading') }}</p>

    <p v-else-if="load_failed" class="py-4 text-sm text-red-400">
      {{ $t('gallery.editor.failedLoad') }}
    </p>

    <template v-else>
      <p
        v-if="!collections.length && !directories.length"
        class="py-4 text-center text-sm text-foreground/40"
      >
        {{ $t('gallery.editor.empty') }}
      </p>

      <section class="flex flex-col gap-2">
        <h2 class="text-lg font-bold">{{ $t('gallery.editor.collections') }}</h2>

        <p v-if="!collections.length" class="py-4 text-center text-sm text-foreground/40">
          {{ $t('gallery.editor.noCollections') }}
        </p>

        <div
          v-for="collection in collections"
          :key="collection.name"
          class="flex items-center gap-3 rounded-lg border border-accent/20 bg-background p-3"
        >
          <CDNImage
            v-if="coverUri(collection)"
            class="h-12 w-12 shrink-0 rounded object-cover"
            :uri="coverUri(collection)!"
            :resolution="120"
            :quality="70"
          ></CDNImage>
          <div
            v-else
            class="flex h-12 w-12 shrink-0 items-center justify-center rounded bg-shade text-foreground/40"
          >
            <Icon icon="mynaui:image" />
          </div>

          <div class="flex min-w-0 flex-1 flex-col">
            <span class="truncate text-sm font-semibold">
              {{ collection.title || collection.name }}
            </span>
            <span class="truncate text-xs text-foreground/50">{{ collection.name }}</span>
            <span class="text-xs text-foreground/40">
              {{ $t('gallery.editor.assetCount', { count: collection.count }) }}
            </span>
          </div>

          <RouterLink
            :to="collectionRoute(collection.name)"
            class="shrink-0 rounded border border-accent/40 px-2 py-1 text-xs hover:bg-accent/10"
          >
            {{ $t('common.edit') }}
          </RouterLink>
        </div>
      </section>

      <section class="flex flex-col gap-2">
        <h2 class="text-lg font-bold">{{ $t('gallery.editor.available') }}</h2>
        <p class="text-sm text-secondary">{{ $t('gallery.editor.availableHint') }}</p>

        <p v-if="!directories.length" class="py-4 text-center text-sm text-foreground/40">
          {{ $t('gallery.editor.noDirectories') }}
        </p>

        <div
          v-for="directory in directories"
          :key="directory.name"
          class="flex items-center justify-between gap-3 rounded-lg border border-dotted border-accent/20 p-3"
        >
          <div class="flex min-w-0 flex-col">
            <span class="truncate text-sm font-semibold">{{ directory.name }}</span>
            <span class="text-xs text-foreground/40">
              {{ $t('gallery.editor.assetCount', { count: directory.count }) }}
            </span>
          </div>

          <button
            class="shrink-0 rounded bg-accent px-3 py-1 text-xs font-semibold text-background hover:brightness-110 disabled:opacity-50"
            :disabled="!!publishing"
            @click="publish_directory(directory)"
          >
            {{
              publishing === directory.name
                ? $t('gallery.editor.publishing')
                : $t('gallery.editor.publish')
            }}
          </button>
        </div>
      </section>
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useToast } from '@/composables/useToast'
import {
  galleryDetailPath,
  getGalleryCollectionRaw,
  getGalleryDirectoriesRaw,
  upsertGalleryCollectionRaw
} from '@/services/gallery'

import type { GalleryCollection, GalleryDirectory } from '@/interfaces/Gallery'

interface CollectionEntry {
  name: string
  title: string
  cover: string
  count: number
}

const { t } = useI18n()
const toast = useToast()

const collections = ref<CollectionEntry[]>([])
const directories = ref<GalleryDirectory[]>([])
const loading = ref(true)
const load_failed = ref(false)
const publishing = ref('')

function coverUri(collection: CollectionEntry): string {
  if (!collection.cover) return ''

  return `${galleryDetailPath(collection.name)}/${collection.cover}`
}

function collectionRoute(name: string): string {
  return `/editor/gallery/${encodeURIComponent(name)}`
}

function isCollectionMap(value: unknown): value is GalleryCollection {
  return !!value && typeof value === 'object' && !Array.isArray(value)
}

async function load(initial = false): Promise<void> {
  if (initial) loading.value = true

  try {
    const [collection_response, directory_response] = await Promise.all([
      getGalleryCollectionRaw().send(),
      getGalleryDirectoriesRaw().send()
    ])

    const collection_map = isCollectionMap(collection_response) ? collection_response : {}
    const directory_list = Array.isArray(directory_response) ? directory_response : []
    const counts = new Map(directory_list.map((directory) => [directory.name, directory.count]))

    collections.value = Object.entries(collection_map).map(([name, detail]) => ({
      name,
      title: detail?.title ?? '',
      cover: detail?.cover ?? '',
      count: counts.get(name) ?? 0
    }))
    directories.value = directory_list.filter((directory) => !(directory.name in collection_map))
    load_failed.value = false
  } catch {
    load_failed.value = true
    toast.error(t('gallery.editor.failedLoad'))
  } finally {
    loading.value = false
  }
}

async function publish_directory(directory: GalleryDirectory): Promise<void> {
  if (publishing.value) return

  publishing.value = directory.name

  try {
    await upsertGalleryCollectionRaw({
      path: directory.name,
      title: directory.name,
      cover: ''
    })
    await load()
  } catch {
    toast.error(t('gallery.editor.publishFailed'))
  } finally {
    publishing.value = ''
  }
}

onMounted(() => load(true))
</script>
