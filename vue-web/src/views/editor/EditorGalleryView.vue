<template>
  <EditorPage>
    <EditorPageHeader
      :title="collection?.title || $t('editor.gallery')"
      :back-to="GALLERY_ROUTE"
      :back-label="$t('editor.ui.backTo', { section: $t('editor.gallery') })"
    >
      <template #actions>
        <template v-if="collection">
          <EditorButton
            variant="primary"
            :disabled="saving || !has_changes || !title.trim()"
            :loading="saving"
            @click="save_all"
          >
            {{ $t('common.save') }}
          </EditorButton>
          <EditorButton
            variant="destructive"
            :disabled="deleting"
            @click="unpublish_collection"
          >
            {{ $t('common.unpublish') }}
          </EditorButton>
        </template>
      </template>
    </EditorPageHeader>

    <EditorLoading v-if="loading"></EditorLoading>

    <EditorError v-else-if="load_failed" :message="t('gallery.editor.failedLoad')">
      <template #actions>
        <EditorButton variant="secondary" @click="load(true)">
          {{ $t('common.retry') }}
        </EditorButton>
      </template>
    </EditorError>

    <EditorError v-else-if="not_found" :message="$t('editor.notFound')">
      <template #actions>
        <EditorButton variant="secondary" :to="GALLERY_ROUTE">
          {{ $t('editor.ui.backTo', { section: $t('editor.gallery') }) }}
        </EditorButton>
      </template>
    </EditorError>

    <div v-else-if="collection" class="flex flex-col gap-6">
      <GalleryZoom v-if="zoom_uri" :uri="zoom_uri" @close="zoom_uri = ''"></GalleryZoom>

      <section class="overflow-hidden rounded-xl border border-accent/20 bg-background">
        <div class="relative min-h-32 w-full md:min-h-48 lg:min-h-96">
          <CDNImage
            v-if="cover"
            class="max-h-32 w-full object-cover object-center brightness-50 blur-[2px] md:max-h-48 lg:max-h-96"
            :uri="cover_uri"
            :resolution="1024"
            :quality="75"
          ></CDNImage>
          <div v-else class="h-full min-h-32 w-full bg-shade md:min-h-48 lg:min-h-96"></div>

          <div
            class="pointer-events-none absolute inset-0 flex flex-col items-center justify-center gap-1 px-4"
          >
            <input
              v-model="title"
              type="text"
              :placeholder="$t('gallery.editor.titlePlaceholder')"
              :aria-label="$t('gallery.editor.titleLabel')"
              class="pointer-events-auto w-full rounded bg-transparent text-center font-serif text-2xl font-bold uppercase tracking-widest text-white outline-none placeholder:text-white/40 focus:bg-black/20 md:text-4xl"
            />
            <p class="text-sm font-medium text-white/80">
              {{ $t('gallery.editor.directory') }}: {{ collection.path }}
            </p>
          </div>

          <div
            class="absolute inset-x-0 bottom-0 flex justify-end gap-2 bg-gradient-to-t from-background/80 to-transparent p-3"
          >
            <EditorButton
              variant="secondary"
              size="sm"
              class="bg-background/80 backdrop-blur"
              @click="picker = true"
            >
              {{ $t('gallery.editor.chooseCover') }}
            </EditorButton>
            <EditorButton
              v-if="cover"
              variant="secondary"
              size="sm"
              class="bg-background/80 backdrop-blur"
              @click="cover = ''"
            >
              {{ $t('gallery.editor.clearCover') }}
            </EditorButton>
          </div>
        </div>
      </section>

      <EditorSection
        :title="$t('gallery.editor.defaults')"
        :description="$t('gallery.editor.defaultsHint')"
      >
        <div class="flex flex-col gap-4">
          <EditorField
            :label="$t('gallery.editor.description')"
            html-for="gallery-default-description"
          >
            <EditorTextarea
              id="gallery-default-description"
              :model-value="values_for('*').description"
              @update:model-value="setDescription('*', $event)"
            ></EditorTextarea>
          </EditorField>

          <EditorField
            :label="$t('gallery.editor.timezoneAdjustment')"
            html-for="gallery-default-timezone"
            class="max-w-xs"
          >
            <EditorSelect
              id="gallery-default-timezone"
              :model-value="values_for('*').tz_adjustment"
              @update:model-value="setTimezone('*', $event)"
            >
              <option value="">{{ $t('gallery.editor.noTimezoneAdjustment') }}</option>
              <option
                v-for="option in timezone_options_for(values_for('*').tz_adjustment)"
                :key="option.value"
                :value="option.value"
              >
                {{ option.label }}
              </option>
            </EditorSelect>
          </EditorField>

          <p class="text-sm font-medium leading-relaxed text-foreground/60">
            {{ $t('gallery.editor.timezoneHint') }}
          </p>
          <p class="text-sm font-medium leading-relaxed text-foreground/60">
            {{ $t('gallery.editor.directoryFixed') }}
          </p>
        </div>
      </EditorSection>

      <GalleryPicker
        v-if="picker"
        :collection="collection.path"
        @select="onCoverSelect"
        @close="picker = false"
      ></GalleryPicker>

      <EditorSection :title="$t('gallery.editor.assets')">
        <EditorEmptyState
          v-if="!assets.length"
          :message="$t('gallery.editor.noAssets')"
        ></EditorEmptyState>

        <div v-else class="grid gap-2 xl:grid-cols-2">
          <div
            v-for="asset in page_assets"
            :key="asset.filename"
            class="flex flex-col gap-3 rounded-lg border border-accent/20 p-3"
          >
            <div class="flex items-start gap-3">
              <CDNImage
                class="h-20 w-20 shrink-0 rounded object-cover lg:cursor-zoom-in"
                :uri="assetUri(asset)"
                :resolution="160"
                :quality="70"
                @click="zoom_uri = assetUri(asset)"
              ></CDNImage>

              <div class="flex min-w-0 flex-1 flex-col gap-1">
                <div class="flex flex-wrap items-center gap-x-2 gap-y-1">
                  <span class="truncate text-lg font-medium" :title="asset.filename">
                    {{ asset.filename }}
                  </span>
                  <span
                    v-if="!hasOverride(asset.filename)"
                    class="rounded-full bg-accent/10 px-2 py-0.5 text-sm font-medium text-foreground/60"
                  >
                    {{ $t('gallery.editor.inherited') }}
                  </span>
                  <span class="text-sm font-medium text-foreground/60">
                    {{ $t('gallery.editor.localTime') }}:
                    {{ localTime(asset) || $t('common.notAvailable') }}
                  </span>
                </div>
              </div>
            </div>

            <EditorTextarea
              :aria-label="$t('gallery.editor.description')"
              :placeholder="$t('gallery.editor.description')"
              :rows="2"
              :model-value="values_for(asset.filename).description"
              @update:model-value="setDescription(asset.filename, $event)"
            ></EditorTextarea>

            <EditorField
              :label="$t('gallery.editor.timezoneAdjustment')"
              :html-for="`gallery-asset-timezone-${asset.filename}`"
              class="w-full sm:max-w-xs"
            >
              <EditorSelect
                :id="`gallery-asset-timezone-${asset.filename}`"
                :model-value="values_for(asset.filename).tz_adjustment"
                @update:model-value="setTimezone(asset.filename, $event)"
              >
                <option value="">{{ $t('gallery.editor.noTimezoneAdjustment') }}</option>
                <option
                  v-for="option in timezone_options_for(
                    values_for(asset.filename).tz_adjustment
                  )"
                  :key="option.value"
                  :value="option.value"
                >
                  {{ option.label }}
                </option>
              </EditorSelect>
            </EditorField>
          </div>
        </div>

        <div v-if="total_pages > 1" class="mt-4 flex items-center justify-between gap-3">
          <EditorButton
            size="sm"
            variant="secondary"
            :disabled="page <= 1"
            @click="page--"
          >
            {{ $t('common.previous') }}
          </EditorButton>
          <span class="text-sm font-medium text-foreground/60">{{ page }} / {{ total_pages }}</span>
          <EditorButton
            size="sm"
            variant="secondary"
            :disabled="page >= total_pages"
            @click="page++"
          >
            {{ $t('common.next') }}
          </EditorButton>
        </div>
      </EditorSection>

      <EditorSection
        v-if="stale_details.length"
        :title="$t('gallery.editor.stale')"
        :description="$t('gallery.editor.staleHint')"
      >
        <div class="flex flex-col gap-2">
          <div
            v-for="detail in stale_details"
            :key="detail.filename"
            class="flex items-center justify-between gap-3 rounded-lg border border-dotted border-accent/20 p-3"
          >
            <span class="truncate text-lg font-medium" :title="detail.filename">
              {{ detail.filename }}
            </span>
            <EditorButton
              size="sm"
              variant="destructive"
              :disabled="deleting_stale === detail.filename"
              @click="removeStale(detail)"
            >
              {{ $t('common.delete') }}
            </EditorButton>
          </div>
        </div>
      </EditorSection>

      <div
        class="sticky bottom-0 flex items-center justify-between gap-3 border-t border-accent/20 bg-background py-3"
      >
        <span class="text-sm font-medium text-foreground/60">
          {{
            $t(
              'gallery.editor.pending',
              { count: pending_count },
              pending_count
            )
          }}
        </span>
      </div>
    </div>
  </EditorPage>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import dayjs from 'dayjs'
import LocalizedFormat from 'dayjs/plugin/localizedFormat'
dayjs.extend(LocalizedFormat)

import GalleryPicker from '@/components/editor/article/GalleryPicker.vue'
import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorEmptyState from '@/components/editor/ui/EditorEmptyState.vue'
import EditorError from '@/components/editor/ui/EditorError.vue'
import EditorField from '@/components/editor/ui/EditorField.vue'
import EditorLoading from '@/components/editor/ui/EditorLoading.vue'
import EditorPage from '@/components/editor/ui/EditorPage.vue'
import EditorPageHeader from '@/components/editor/ui/EditorPageHeader.vue'
import EditorSection from '@/components/editor/ui/EditorSection.vue'
import EditorSelect from '@/components/editor/ui/EditorSelect.vue'
import EditorTextarea from '@/components/editor/ui/EditorTextarea.vue'
import CDNImage from '@/components/shared/CDNImage.vue'
import { useSaveShortcut } from '@/composables/useSaveShortcut'
import { useToast } from '@/composables/useToast'
import { useUnsavedGuard } from '@/composables/useUnsavedGuard'
import {
  deleteGalleryCollectionRaw,
  deleteGalleryDetailRaw,
  galleryDetailPath,
  getGalleryCollectionRaw,
  getGalleryDetailsRaw,
  loadMoreAssets,
  upsertGalleryCollectionRaw,
  upsertGalleryDetailRaw
} from '@/services/gallery'

import type {
  Asset,
  GalleryCollection,
  GalleryCollectionDetail,
  GalleryDetail
} from '@/interfaces/Gallery'

interface OverrideValues {
  description: string
  tz_adjustment: string
}

const PAGE_SIZE = 40
const GALLERY_ROUTE = '/editor/gallery'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const collection_name = computed(() =>
  typeof route.params.collection === 'string' ? route.params.collection : ''
)

const loading = ref(true)
const load_failed = ref(false)
const not_found = ref(false)
const leaving = ref(false)
const collection = ref<GalleryCollectionDetail | null>(null)

const title = ref('')
const cover = ref('')
const picker = ref(false)
const zoom_uri = ref('')
const saving_metadata = ref(false)
const deleting = ref(false)

const assets = ref<Asset[]>([])
const details = ref<GalleryDetail[]>([])
const pending = ref<Record<string, OverrideValues>>({})
const saving_overrides = ref(false)
const deleting_stale = ref('')
const page = ref(1)

const wildcard_row = computed(() => details.value.find((detail) => detail.filename === '*'))

const cover_uri = computed(() => {
  if (!cover.value) return ''

  return `${galleryDetailPath(collection_name.value)}/${cover.value}`
})

const total_pages = computed(() => Math.max(1, Math.ceil(assets.value.length / PAGE_SIZE)))

const page_assets = computed(() =>
  assets.value.slice((page.value - 1) * PAGE_SIZE, page.value * PAGE_SIZE)
)

const asset_filenames = computed(() => new Set(assets.value.map((asset) => asset.filename)))

const stale_details = computed(() =>
  details.value.filter(
    (detail) => detail.filename !== '*' && !asset_filenames.value.has(detail.filename)
  )
)

const pending_count = computed(() => Object.keys(pending.value).length)

const metadata_dirty = computed(
  () =>
    title.value !== (collection.value?.title ?? '') ||
    cover.value !== (collection.value?.cover ?? '')
)

const has_changes = computed(() => metadata_dirty.value || pending_count.value > 0)
const saving = computed(() => saving_metadata.value || saving_overrides.value)

const dirty = computed(() => {
  if (leaving.value) return false
  if (metadata_dirty.value) return true

  return Object.entries(pending.value).some(([filename, values]) => {
    const inherited = inherited_values(filename)
    return (
      values.description !== inherited.description ||
      values.tz_adjustment !== inherited.tz_adjustment
    )
  })
})

useUnsavedGuard(dirty)
useSaveShortcut(save_all, { enabled: () => has_changes.value && !saving.value })

function isCollectionMap(value: unknown): value is GalleryCollection {
  return !!value && typeof value === 'object' && !Array.isArray(value)
}

function row_for(filename: string): GalleryDetail | undefined {
  if (filename === '*') return wildcard_row.value

  return details.value.find((detail) => detail.filename === filename)
}

function inherited_values(filename: string): OverrideValues {
  const source = row_for(filename) ?? wildcard_row.value

  return {
    description: source?.description ?? '',
    tz_adjustment: source?.tz_adjustment ? String(source.tz_adjustment) : ''
  }
}

function values_for(filename: string): OverrideValues {
  return pending.value[filename] ?? inherited_values(filename)
}

const BASE_TIMEZONE_OFFSETS = Array.from({ length: 27 }, (_, index) => index - 12)

function timezone_label(offset: number): string {
  const sign = offset >= 0 ? '+' : '-'

  return `UTC${sign}${String(Math.abs(offset)).padStart(2, '0')}:00`
}

function timezone_options_for(value: string): { value: string; label: string }[] {
  const offsets = [...BASE_TIMEZONE_OFFSETS]

  if (value !== '') {
    const parsed = Number(value)

    if (Number.isFinite(parsed) && !offsets.includes(parsed)) {
      offsets.push(parsed)
      offsets.sort((left, right) => left - right)
    }
  }

  return offsets.map((offset) => ({ value: String(offset), label: timezone_label(offset) }))
}

function update(filename: string, field: keyof OverrideValues, value: string): void {
  const next = { ...values_for(filename), [field]: value }
  const inherited = inherited_values(filename)
  const next_pending = { ...pending.value }

  if (
    next.description === inherited.description &&
    next.tz_adjustment === inherited.tz_adjustment
  ) {
    delete next_pending[filename]
  } else {
    next_pending[filename] = next
  }

  pending.value = next_pending
}

function setDescription(filename: string, value: string | number | null): void {
  update(filename, 'description', value == null ? '' : String(value))
}

function setTimezone(filename: string, value: string | number | null): void {
  update(filename, 'tz_adjustment', value == null ? '' : String(value))
}

function hasOverride(filename: string): boolean {
  return details.value.some((detail) => detail.filename === filename)
}

function assetUri(asset: Asset): string {
  const path = asset.collection || galleryDetailPath(collection_name.value)

  return `${path}/${asset.filename}`
}

function localTime(asset: Asset): string {
  if (!asset.exif?.datetime) return ''

  const base = dayjs(asset.exif.datetime)
  const adjustment = values_for(asset.filename).tz_adjustment
  const hours = Number(adjustment)

  if (adjustment === '' || !Number.isFinite(hours) || hours === 0) {
    return base.format('LLL')
  }

  return base.add(hours, 'h').format('LLL')
}

async function loadAssets(name: string): Promise<Asset[]> {
  const first = await loadMoreAssets(name, 1).send()
  const loaded = Array.isArray(first?.data) ? [...first.data] : []
  const pages = Number(first?.total_pages ?? 0)

  for (let next = 2; next <= pages; next++) {
    const response = await loadMoreAssets(name, next).send()

    if (Array.isArray(response?.data)) loaded.push(...response.data)
  }

  return loaded
}

async function loadDetails(name: string): Promise<GalleryDetail[]> {
  const loaded: GalleryDetail[] = []
  let current = 1
  let pages = 1

  do {
    const response = await getGalleryDetailsRaw(name, current).send()

    if (Array.isArray(response?.data)) loaded.push(...response.data)
    pages = Number(response?.total_pages ?? 0)
    current++
  } while (current <= pages)

  return loaded
}

async function load(initial = false): Promise<void> {
  const name = collection_name.value

  if (initial) loading.value = true

  load_failed.value = false
  not_found.value = false
  collection.value = null
  assets.value = []
  details.value = []
  pending.value = {}
  page.value = 1
  zoom_uri.value = ''

  if (!name) {
    not_found.value = true
    loading.value = false
    return
  }

  try {
    const response = await getGalleryCollectionRaw().send()
    const map = isCollectionMap(response) ? response : {}
    const entry = map[name]

    if (!entry) {
      not_found.value = true
      return
    }

    collection.value = {
      path: name,
      title: entry.title ?? '',
      cover: entry.cover ?? ''
    }
    title.value = collection.value.title ?? ''
    cover.value = collection.value.cover ?? ''

    const [loaded_assets, loaded_details] = await Promise.all([loadAssets(name), loadDetails(name)])

    assets.value = loaded_assets
    details.value = loaded_details
  } catch {
    load_failed.value = true
  } finally {
    loading.value = false
  }
}

function onCoverSelect(selection: { collection: string; filename: string }): void {
  cover.value = selection.filename
  picker.value = false
}

async function saveMetadata(): Promise<boolean> {
  const name = collection_name.value
  const next_title = title.value.trim()

  if (!collection.value || !next_title || saving_metadata.value) return false

  saving_metadata.value = true

  try {
    const saved = await upsertGalleryCollectionRaw({
      path: name,
      title: next_title,
      cover: cover.value
    })

    collection.value = {
      path: name,
      title: saved?.title ?? next_title,
      cover: saved?.cover ?? cover.value
    }
    title.value = collection.value.title ?? ''
    cover.value = collection.value.cover ?? ''
    return true
  } catch {
    toast.error(t('gallery.editor.failedSave'))
    return false
  } finally {
    saving_metadata.value = false
  }
}

async function unpublish_collection(): Promise<void> {
  const name = collection_name.value

  if (
    !name ||
    deleting.value ||
    !confirm(t('gallery.editor.confirmUnpublish', { title: title.value || name }))
  ) {
    return
  }

  deleting.value = true

  try {
    await deleteGalleryCollectionRaw(name)
    leaving.value = true
    pending.value = {}
    toast.success(t('common.done'))
    await router.push(GALLERY_ROUTE)
  } catch {
    toast.error(t('gallery.editor.failedUnpublish'))
  } finally {
    deleting.value = false
  }
}

async function saveOverrides(): Promise<boolean> {
  const name = collection_name.value

  if (!name || saving_overrides.value || !pending_count.value) return true

  saving_overrides.value = true

  const saved: string[] = []
  let failed = false

  for (const [filename, values] of Object.entries(pending.value)) {
    const description = values.description.trim()
    const raw_adjustment = values.tz_adjustment.trim()
    const adjustment = raw_adjustment === '' ? 0 : Number(raw_adjustment)

    if (raw_adjustment !== '' && !Number.isFinite(adjustment)) {
      failed = true
      continue
    }

    try {
      if (!description && raw_adjustment === '') {
        if (hasOverride(filename)) await deleteGalleryDetailRaw(name, filename)
      } else {
        await upsertGalleryDetailRaw(name, {
          path: galleryDetailPath(name),
          filename,
          description,
          tz_adjustment: adjustment
        })
      }

      saved.push(filename)
    } catch {
      failed = true
    }
  }

  if (saved.length) {
    const next_pending = { ...pending.value }
    for (const filename of saved) delete next_pending[filename]
    pending.value = next_pending

    try {
      details.value = await loadDetails(name)
    } catch {
      failed = true
    }
  }

  if (failed) toast.error(t('gallery.editor.failedSave'))
  saving_overrides.value = false

  return !failed
}

async function removeStale(detail: GalleryDetail): Promise<void> {
  const name = collection_name.value

  if (
    !name ||
    deleting_stale.value ||
    !confirm(t('gallery.editor.confirmDeleteStale', { filename: detail.filename }))
  ) {
    return
  }

  deleting_stale.value = detail.filename

  try {
    await deleteGalleryDetailRaw(name, detail.filename)
    details.value = details.value.filter((row) => row.filename !== detail.filename)
    toast.success(t('common.deleted'))
  } catch {
    toast.error(t('gallery.editor.failedDelete'))
  } finally {
    deleting_stale.value = ''
  }
}

watch(collection_name, () => load(true), { immediate: true })

async function save_all(): Promise<void> {
  if (saving.value) return

  if (!has_changes.value) return

  if (metadata_dirty.value && !title.value.trim()) {
    toast.error(t('gallery.editor.titleRequired'))
    return
  }

  let ok = true

  if (metadata_dirty.value) ok = await saveMetadata()
  if (ok && pending_count.value) ok = await saveOverrides()

  if (ok) toast.success(t('common.done'))
}
</script>
