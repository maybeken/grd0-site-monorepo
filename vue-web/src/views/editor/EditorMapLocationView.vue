<template>
  <EditorPage>
    <EditorPageHeader
      :title="form.title || $t('map.editor.editTitle')"
      :back-to="MAP_ROUTE"
      :back-label="$t('editor.ui.backTo', { section: $t('editor.mapLocations') })"
    >
      <template #actions>
        <template v-if="location">
          <EditorButton variant="secondary" to="/travel/map" target="_blank">
            {{ $t('map.editor.publicMap') }}
          </EditorButton>
          <EditorButton variant="secondary" @click="form.published = !form.published">
            {{ form.published ? $t('common.unpublish') : $t('common.publish') }}
          </EditorButton>
          <EditorButton
            variant="primary"
            :disabled="!can_save"
            :loading="saving"
            @click="save"
          >
            {{ $t('common.save') }}
          </EditorButton>
        </template>
      </template>
    </EditorPageHeader>

    <EditorLoading v-if="loading"></EditorLoading>

    <EditorError v-else-if="load_failed" :message="t('map.editor.failedLoad')">
      <template #actions>
        <EditorButton variant="secondary" @click="load">
          {{ $t('common.retry') }}
        </EditorButton>
      </template>
    </EditorError>

    <EditorError v-else-if="!location" :message="$t('editor.notFound')">
      <template #actions>
        <EditorButton variant="secondary" :to="MAP_ROUTE">
          {{ $t('editor.ui.backTo', { section: $t('editor.mapLocations') }) }}
        </EditorButton>
      </template>
    </EditorError>

    <div v-else class="flex flex-col gap-6">
      <div class="flex flex-wrap items-center gap-x-3 gap-y-2">
        <span class="font-mono text-sm font-medium text-foreground/60">
          {{ $t('map.editor.slugField') }}: {{ location.slug }}
        </span>
        <span class="text-sm font-medium text-foreground/60">
          {{ $t('map.editor.slugImmutable') }}
        </span>
        <span
          class="rounded-full px-2 py-0.5 text-sm font-medium uppercase tracking-wide"
          :class="
            form.published
              ? 'bg-emerald-600/20 text-emerald-400'
              : 'bg-orange-500/20 text-orange-400'
          "
        >
          {{ form.published ? $t('map.editor.published') : $t('map.editor.draft') }}
        </span>
        <span class="text-sm font-medium text-foreground/60">
          {{ form.published ? $t('map.editor.publishedHint') : $t('map.editor.draftHint') }}
        </span>
      </div>

      <div class="grid gap-6 lg:grid-cols-2">
        <div class="flex flex-col gap-6">
          <EditorSection :title="$t('map.editor.details')">
            <div class="flex flex-col gap-4">
              <EditorField
                :label="$t('map.editor.titleField')"
                html-for="map-location-title"
                required
                :error="title_error ? $t('map.editor.titleRequired') : ''"
              >
                <EditorInput
                  id="map-location-title"
                  v-model="form.title"
                  :placeholder="$t('map.editor.titlePlaceholder')"
                  required
                ></EditorInput>
              </EditorField>

              <EditorField
                :label="$t('map.editor.subtitleField')"
                html-for="map-location-subtitle"
              >
                <EditorInput
                  id="map-location-subtitle"
                  v-model="form.subtitle"
                  :placeholder="$t('map.editor.subtitlePlaceholder')"
                ></EditorInput>
              </EditorField>

              <div class="flex flex-col gap-1.5">
                <p class="text-lg font-medium">{{ $t('map.editor.iconField') }}</p>
                <div class="flex items-center gap-2">
                  <IconSelector v-model="form.icon" class="min-w-0 flex-1" />
                  <Icon v-if="form.icon" :icon="form.icon" height="auto" />
                </div>
              </div>

              <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <EditorField
                  :label="$t('map.editor.colorField')"
                  html-for="map-location-color"
                >
                  <EditorSelect id="map-location-color" v-model="form.color">
                    <option v-for="key in COLOR_KEYS" :key="key" :value="key">
                      {{ $t(`map.editor.color.${key}`) }}
                    </option>
                  </EditorSelect>
                </EditorField>

                <EditorField
                  :label="$t('map.editor.textColorField')"
                  html-for="map-location-text-color"
                >
                  <EditorSelect id="map-location-text-color" v-model="form.text_color">
                    <option value="">{{ $t('map.editor.textColor.default') }}</option>
                    <option v-for="key in TEXT_COLOR_KEYS" :key="key" :value="key">
                      {{ $t(`map.editor.textColor.${key}`) }}
                    </option>
                  </EditorSelect>
                </EditorField>
              </div>
            </div>
          </EditorSection>

          <EditorSection :title="$t('map.editor.visibility')">
            <div class="flex flex-col gap-4">
              <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <EditorField
                  :label="$t('map.editor.visibleFrom')"
                  html-for="map-location-visible-from"
                >
                  <EditorInput
                    id="map-location-visible-from"
                    v-model="form.display_zoom"
                    type="number"
                    step="0.01"
                    :placeholder="$t('map.editor.zoomPlaceholder')"
                  ></EditorInput>
                </EditorField>

                <EditorField
                  :label="$t('map.editor.visibleUntil')"
                  html-for="map-location-visible-until"
                >
                  <EditorInput
                    id="map-location-visible-until"
                    v-model="form.hide_zoom"
                    type="number"
                    step="0.01"
                    :placeholder="$t('map.editor.zoomPlaceholder')"
                  ></EditorInput>
                </EditorField>
              </div>

              <p class="text-sm font-medium leading-relaxed text-foreground/60">
                {{ $t('map.editor.visibilityHint') }}
              </p>
            </div>
          </EditorSection>

          <EditorSection :title="$t('map.editor.coordinates')">
            <div class="flex flex-col gap-4">
              <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <EditorField
                  :label="$t('map.editor.latitude')"
                  html-for="map-location-latitude"
                  :error="latitude_error ? $t('map.editor.invalidLatitude') : ''"
                >
                  <EditorInput
                    id="map-location-latitude"
                    v-model="form.latitude"
                    type="number"
                    step="any"
                    :invalid="latitude_error"
                  ></EditorInput>
                </EditorField>

                <EditorField
                  :label="$t('map.editor.longitude')"
                  html-for="map-location-longitude"
                  :error="longitude_error ? $t('map.editor.invalidLongitude') : ''"
                >
                  <EditorInput
                    id="map-location-longitude"
                    v-model="form.longitude"
                    type="number"
                    step="any"
                    :invalid="longitude_error"
                  ></EditorInput>
                </EditorField>
              </div>

              <p class="text-sm font-medium leading-relaxed text-foreground/60">
                {{ $t('map.editor.coordinateHint') }}
              </p>
            </div>
          </EditorSection>
        </div>

        <div class="flex flex-col gap-6">
          <EditorSection :title="$t('map.editor.preview')">
            <div class="flex min-h-24 items-center justify-center rounded-lg bg-shade p-6">
              <div
                class="max-w-72 select-none rounded-r-2xl rounded-b-2xl px-4 py-2"
                :class="[
                  parseTailwindColor('bg', form.color),
                  parseTailwindColor('text', form.text_color || 'default')
                ]"
              >
                <div class="flex gap-1 pb-1">
                  <Icon
                    v-if="form.icon"
                    :class="form.title ? ['mr-1'] : []"
                    :icon="form.icon"
                    height="auto"
                  />
                  <span class="text-lg font-medium">{{ form.title || '' }}</span>
                </div>
                <p v-if="form.subtitle" class="text-sm italic text-justify">
                  {{ form.subtitle }}
                </p>
              </div>
            </div>
          </EditorSection>

          <EditorSection :title="$t('map.editor.coordinates')">
            <LocationMapPicker
              :key="slug"
              :latitude="latitude"
              :longitude="longitude"
              :title="form.title"
              :subtitle="form.subtitle"
              :icon="form.icon"
              :color="form.color"
              :text_color="form.text_color"
              @select="onSelectCoordinates"
            />
          </EditorSection>
        </div>
      </div>
    </div>
  </EditorPage>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorError from '@/components/editor/ui/EditorError.vue'
import EditorField from '@/components/editor/ui/EditorField.vue'
import EditorInput from '@/components/editor/ui/EditorInput.vue'
import EditorLoading from '@/components/editor/ui/EditorLoading.vue'
import EditorPage from '@/components/editor/ui/EditorPage.vue'
import EditorPageHeader from '@/components/editor/ui/EditorPageHeader.vue'
import EditorSection from '@/components/editor/ui/EditorSection.vue'
import EditorSelect from '@/components/editor/ui/EditorSelect.vue'
import LocationMapPicker from '@/components/map/LocationMapPicker.vue'
import IconSelector from '@/components/editor/IconSelector.vue'
import Icon from '@/components/shared/Icon.vue'
import { useSaveShortcut } from '@/composables/useSaveShortcut'
import { useToast } from '@/composables/useToast'
import { useUnsavedGuard } from '@/composables/useUnsavedGuard'
import {
  formatZoom,
  isValidLatitude,
  isValidLongitude,
  parseZoom,
  zoomToResolution
} from '@/helpers/map'
import { COLOR_KEYS, TEXT_COLOR_KEYS, parseTailwindColor } from '@/helpers/mapPalette'
import { getMapLocationBySlug, upsertMapLocation } from '@/services/travelersMap'

import type { MapLocation } from '@/interfaces/TravelersMap'

const MAP_ROUTE = '/editor/travel'

const route = useRoute()
const { t } = useI18n()
const toast = useToast()

const location = ref<MapLocation | null>(null)
const loading = ref(true)
const load_failed = ref(false)
const saving = ref(false)

const form = reactive({
  title: '',
  subtitle: '',
  icon: '',
  color: 'default',
  text_color: '',
  published: false,
  latitude: '' as string | number,
  longitude: '' as string | number,
  display_zoom: '' as string | number,
  hide_zoom: '' as string | number
})

const baseline = ref('')
const baseline_display_zoom = ref<string | number>('')
const baseline_hide_zoom = ref<string | number>('')
const original_display_at = ref(0)
const original_hide_at = ref(0)

const slug = computed(() => {
  const value = route.params.slug

  return typeof value === 'string' ? value : ''
})

function snapshot(): string {
  return JSON.stringify(form)
}

const dirty = computed(() => snapshot() !== baseline.value)
const title_error = computed(() => form.title.trim() === '')

useUnsavedGuard(dirty)
useSaveShortcut(save, { enabled: () => can_save.value })

function toNumber(value: string | number): number | null {
  if (typeof value === 'number') {
    return Number.isFinite(value) ? value : null
  }

  if (value.trim() === '') {
    return null
  }

  const number = Number(value)

  return Number.isFinite(number) ? number : null
}

const latitude = computed(() => toNumber(form.latitude))
const longitude = computed(() => toNumber(form.longitude))

const latitude_error = computed(
  () =>
    String(form.latitude).trim() !== '' &&
    (latitude.value === null || !isValidLatitude(latitude.value))
)
const longitude_error = computed(
  () =>
    String(form.longitude).trim() !== '' &&
    (longitude.value === null || !isValidLongitude(longitude.value))
)

const can_save = computed(
  () =>
    dirty.value &&
    !saving.value &&
    form.title.trim() !== '' &&
    latitude.value !== null &&
    isValidLatitude(latitude.value) &&
    longitude.value !== null &&
    isValidLongitude(longitude.value)
)

function applyLocation(value: MapLocation): void {
  const lat = value.pos?.[0]
  const lon = value.pos?.[1]

  form.title = value.title ?? ''
  form.subtitle = value.subtitle ?? ''
  form.icon = value.icon ?? ''
  form.color = value.color || 'default'
  form.text_color = value.text_color ?? ''
  form.published = Boolean(value.published)
  form.latitude = typeof lat === 'number' ? String(lat) : ''
  form.longitude = typeof lon === 'number' ? String(lon) : ''
  form.display_zoom = formatZoom(value.display_at)
  form.hide_zoom = formatZoom(value.hide_at)

  original_display_at.value = value.display_at ?? 0
  original_hide_at.value = value.hide_at ?? 0
  baseline_display_zoom.value = form.display_zoom
  baseline_hide_zoom.value = form.hide_zoom
  baseline.value = snapshot()
}

async function load(): Promise<void> {
  loading.value = true
  load_failed.value = false
  location.value = null

  try {
    location.value = await getMapLocationBySlug(slug.value)

    if (location.value) {
      applyLocation(location.value)
    }
  } catch {
    load_failed.value = true
  } finally {
    loading.value = false
  }
}

function thresholdFromZoom(value: string | number): number {
  const zoom = parseZoom(value)

  if (zoom === null) {
    return 0
  }

  const resolution = zoomToResolution(zoom)

  return Number.isFinite(resolution) ? Math.round(resolution) : 0
}

function onSelectCoordinates(coordinates: number[]): void {
  const lat = coordinates[0]
  const lon = coordinates[1]

  if (typeof lat !== 'number' || typeof lon !== 'number') {
    return
  }

  form.latitude = lat.toFixed(6)
  form.longitude = lon.toFixed(6)
}

async function save(): Promise<void> {
  const current = location.value
  const lat = latitude.value
  const lon = longitude.value

  if (!current || !can_save.value || lat === null || lon === null) {
    return
  }

  const display_at =
    form.display_zoom === baseline_display_zoom.value
      ? original_display_at.value
      : thresholdFromZoom(form.display_zoom)
  const hide_at =
    form.hide_zoom === baseline_hide_zoom.value
      ? original_hide_at.value
      : thresholdFromZoom(form.hide_zoom)

  saving.value = true

  try {
    const title = form.title.trim()
    const subtitle = form.subtitle.trim()
    const icon = form.icon.trim()

    await upsertMapLocation({
      slug: current.slug,
      title,
      subtitle,
      icon,
      color: form.color,
      text_color: form.text_color,
      pos: [lat, lon],
      display_at,
      hide_at,
      published: form.published
    })

    form.title = title
    form.subtitle = subtitle
    form.icon = icon
    original_display_at.value = display_at
    original_hide_at.value = hide_at
    baseline_display_zoom.value = form.display_zoom
    baseline_hide_zoom.value = form.hide_zoom
    location.value = {
      ...current,
      title,
      subtitle,
      icon,
      color: form.color,
      text_color: form.text_color,
      pos: [lat, lon],
      display_at,
      hide_at,
      published: form.published
    }
    baseline.value = snapshot()
    toast.success(t('map.editor.saved'))
  } catch {
    toast.error(t('map.editor.failedSave'))
  } finally {
    saving.value = false
  }
}

watch(slug, load, { immediate: true })
</script>
