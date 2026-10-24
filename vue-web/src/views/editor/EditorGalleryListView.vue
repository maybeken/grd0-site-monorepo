<template>
  <EditorPage>
    <EditorPageHeader
      :title="$t('editor.gallery')"
      :description="$t('gallery.editor.listHint')"
    ></EditorPageHeader>

    <EditorError v-if="load_failed" :message="t('gallery.editor.failedLoad')">
      <template #actions>
        <EditorButton variant="secondary" @click="load(true)">
          {{ $t('common.retry') }}
        </EditorButton>
      </template>
    </EditorError>

    <EditorLoading v-else-if="loading"></EditorLoading>

    <EditorEmptyState
      v-else-if="!collections.length && !directories.length"
      :message="$t('gallery.editor.empty')"
    ></EditorEmptyState>

    <template v-else>
      <EditorSection :title="$t('gallery.editor.collections')">
        <EditorEmptyState
          v-if="!collections.length"
          :message="$t('gallery.editor.noCollections')"
        ></EditorEmptyState>

        <div v-else class="flex flex-col gap-2">
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
              class="flex h-12 w-12 shrink-0 items-center justify-center rounded bg-shade text-foreground/60"
            >
              <Icon icon="mynaui:image"></Icon>
            </div>

            <div class="flex min-w-0 flex-1 flex-col">
              <span class="truncate text-lg font-medium" :title="collection.title || collection.name">
                {{ collection.title || collection.name }}
              </span>
              <span class="truncate text-sm font-medium text-foreground/60">
                {{ collection.name }}
              </span>
              <span class="text-sm font-medium text-foreground/60">
                {{
                  $t(
                    'gallery.editor.assetCount',
                    { count: collection.count },
                    collection.count
                  )
                }}
              </span>
            </div>

            <EditorButton
              size="sm"
              variant="secondary"
              :to="collectionRoute(collection.name)"
            >
              {{ $t('common.edit') }}
            </EditorButton>
          </div>
        </div>
      </EditorSection>

      <EditorSection
        :title="$t('gallery.editor.available')"
        :description="$t('gallery.editor.availableHint')"
      >
        <EditorEmptyState
          v-if="!directories.length"
          :message="$t('gallery.editor.noDirectories')"
        ></EditorEmptyState>

        <div v-else class="flex flex-col gap-2">
          <div
            v-for="directory in directories"
            :key="directory.name"
            class="flex items-center justify-between gap-3 rounded-lg border border-dotted border-accent/20 p-3"
          >
            <div class="flex min-w-0 flex-col">
              <span class="truncate text-lg font-medium" :title="directory.name">
                {{ directory.name }}
              </span>
              <span class="text-sm font-medium text-foreground/60">
                {{
                  $t(
                    'gallery.editor.assetCount',
                    { count: directory.count },
                    directory.count
                  )
                }}
              </span>
            </div>

            <EditorButton
              size="sm"
              variant="secondary"
              :disabled="!!publishing"
              :loading="publishing === directory.name"
              @click="publish_directory(directory)"
            >
              {{
                publishing === directory.name
                  ? $t('gallery.editor.publishing')
                  : $t('common.publish')
              }}
            </EditorButton>
          </div>
        </div>
      </EditorSection>
    </template>
  </EditorPage>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorEmptyState from '@/components/editor/ui/EditorEmptyState.vue'
import EditorError from '@/components/editor/ui/EditorError.vue'
import EditorLoading from '@/components/editor/ui/EditorLoading.vue'
import EditorPage from '@/components/editor/ui/EditorPage.vue'
import EditorPageHeader from '@/components/editor/ui/EditorPageHeader.vue'
import EditorSection from '@/components/editor/ui/EditorSection.vue'
import CDNImage from '@/components/shared/CDNImage.vue'
import Icon from '@/components/shared/Icon.vue'
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
    toast.success(t('common.done'))
    await load()
  } catch {
    toast.error(t('gallery.editor.publishFailed'))
  } finally {
    publishing.value = ''
  }
}

onMounted(() => load(true))
</script>
