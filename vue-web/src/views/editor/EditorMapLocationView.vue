<template>
  <div class="flex flex-col gap-4">
    <router-link to="/editor/travel" class="text-sm text-secondary hover:underline">
      {{ $t('editor.backToList') }}
    </router-link>

    <p v-if="loading" class="text-secondary">{{ $t('common.loading') }}</p>

    <div v-else-if="location" class="flex max-w-5xl flex-col gap-6">
      <header class="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
        <div class="flex flex-col gap-1">
          <h1 class="text-2xl font-bold">{{ $t('map.editor.editTitle') }}<CursorBlink /></h1>
          <p class="font-mono text-sm text-secondary">
            {{ $t('map.editor.slugField') }}: {{ location.slug }}
          </p>
          <p class="text-xs text-foreground/40">{{ $t('map.editor.slugImmutable') }}</p>
          <p class="flex items-center gap-2 text-xs">
            <span
              class="rounded-full px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wide"
              :class="
                form.published
                  ? 'bg-emerald-600/20 text-emerald-400'
                  : 'bg-orange-500/20 text-orange-400'
              "
            >
              {{ form.published ? $t('map.editor.published') : $t('map.editor.draft') }}
            </span>
            <span class="text-foreground/40">
              {{ form.published ? $t('map.editor.publishedHint') : $t('map.editor.draftHint') }}
            </span>
          </p>
        </div>

        <div class="flex shrink-0 flex-wrap items-center gap-2">
          <RouterLink
            to="/travel/map"
            class="rounded border border-accent/40 px-3 py-1 text-sm hover:bg-accent/10"
          >
            {{ $t('map.editor.publicMap') }}
          </RouterLink>
          <button
            type="button"
            class="rounded border border-accent/40 px-3 py-1 text-sm hover:bg-accent/10"
            @click="form.published = !form.published"
          >
            {{ form.published ? $t('map.editor.unpublish') : $t('map.editor.publish') }}
          </button>
          <button
            class="rounded bg-accent px-3 py-1 text-sm font-semibold text-background hover:brightness-110 disabled:opacity-50"
            :disabled="!can_save"
            @click="save"
          >
            {{ saving ? $t('common.saving') : $t('common.save') }}
          </button>
        </div>
      </header>

      <div class="grid gap-6 lg:grid-cols-2">
        <div class="flex flex-col gap-6">
          <section class="flex flex-col gap-4 rounded-xl border border-accent/20 bg-background p-4">
            <h2 class="text-lg font-bold">{{ $t('map.editor.details') }}</h2>

            <div class="flex flex-col gap-1">
              <label class="text-sm font-semibold" for="map-location-title">
                {{ $t('map.editor.titleField') }}
              </label>
              <input
                id="map-location-title"
                v-model="form.title"
                class="rounded border border-accent/40 bg-background px-2 py-1 text-sm"
                :placeholder="$t('map.editor.titlePlaceholder')"
              />
              <p v-if="!form.title.trim()" class="text-xs text-red-400">
                {{ $t('map.editor.titleRequired') }}
              </p>
            </div>

            <div class="flex flex-col gap-1">
              <label class="text-sm font-semibold" for="map-location-subtitle">
                {{ $t('map.editor.subtitleField') }}
              </label>
              <input
                id="map-location-subtitle"
                v-model="form.subtitle"
                class="rounded border border-accent/40 bg-background px-2 py-1 text-sm"
                :placeholder="$t('map.editor.subtitlePlaceholder')"
              />
            </div>

            <div class="flex flex-col gap-1">
              <label class="text-sm font-semibold">
                {{ $t('map.editor.iconField') }}
              </label>
              <div class="flex items-center gap-2">
                <IconSelector v-model="form.icon" class="min-w-0 flex-1" />
                <Icon v-if="form.icon" :icon="form.icon" height="auto" />
              </div>
            </div>

            <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div class="flex flex-col gap-1">
                <label class="text-sm font-semibold" for="map-location-color">
                  {{ $t('map.editor.colorField') }}
                </label>
                <select
                  id="map-location-color"
                  v-model="form.color"
                  class="rounded border border-accent/40 bg-background px-2 py-1 text-sm"
                >
                  <option v-for="key in COLOR_KEYS" :key="key" :value="key">
                    {{ $t(`map.editor.color.${key}`) }}
                  </option>
                </select>
              </div>

              <div class="flex flex-col gap-1">
                <label class="text-sm font-semibold" for="map-location-text-color">
                  {{ $t('map.editor.textColorField') }}
                </label>
                <select
                  id="map-location-text-color"
                  v-model="form.text_color"
                  class="rounded border border-accent/40 bg-background px-2 py-1 text-sm"
                >
                  <option value="">{{ $t('map.editor.textColor.default') }}</option>
                  <option v-for="key in TEXT_COLOR_KEYS" :key="key" :value="key">
                    {{ $t(`map.editor.textColor.${key}`) }}
                  </option>
                </select>
              </div>
            </div>
          </section>

          <section class="flex flex-col gap-4 rounded-xl border border-accent/20 bg-background p-4">
            <h2 class="text-lg font-bold">{{ $t('map.editor.visibility') }}</h2>

            <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div class="flex flex-col gap-1">
                <label class="text-sm font-semibold" for="map-location-visible-from">
                  {{ $t('map.editor.visibleFrom') }}
                </label>
                <input
                  id="map-location-visible-from"
                  v-model="form.display_zoom"
                  type="number"
                  step="0.01"
                  class="rounded border border-accent/40 bg-background px-2 py-1 text-sm"
                  :placeholder="$t('map.editor.zoomPlaceholder')"
                />
              </div>

              <div class="flex flex-col gap-1">
                <label class="text-sm font-semibold" for="map-location-visible-until">
                  {{ $t('map.editor.visibleUntil') }}
                </label>
                <input
                  id="map-location-visible-until"
                  v-model="form.hide_zoom"
                  type="number"
                  step="0.01"
                  class="rounded border border-accent/40 bg-background px-2 py-1 text-sm"
                  :placeholder="$t('map.editor.zoomPlaceholder')"
                />
              </div>
            </div>

            <p class="text-xs text-foreground/50">{{ $t('map.editor.visibilityHint') }}</p>
          </section>

          <section class="flex flex-col gap-4 rounded-xl border border-accent/20 bg-background p-4">
            <h2 class="text-lg font-bold">{{ $t('map.editor.coordinates') }}</h2>

            <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div class="flex flex-col gap-1">
                <label class="text-sm font-semibold" for="map-location-latitude">
                  {{ $t('map.editor.latitude') }}
                </label>
                <input
                  id="map-location-latitude"
                  v-model="form.latitude"
                  type="number"
                  step="any"
                  class="rounded border border-accent/40 bg-background px-2 py-1 text-sm"
                  :class="latitude_error ? 'border-red-500/60' : ''"
                />
                <p v-if="latitude_error" class="text-xs text-red-400">
                  {{ $t('map.editor.invalidLatitude') }}
                </p>
              </div>

              <div class="flex flex-col gap-1">
                <label class="text-sm font-semibold" for="map-location-longitude">
                  {{ $t('map.editor.longitude') }}
                </label>
                <input
                  id="map-location-longitude"
                  v-model="form.longitude"
                  type="number"
                  step="any"
                  class="rounded border border-accent/40 bg-background px-2 py-1 text-sm"
                  :class="longitude_error ? 'border-red-500/60' : ''"
                />
                <p v-if="longitude_error" class="text-xs text-red-400">
                  {{ $t('map.editor.invalidLongitude') }}
                </p>
              </div>
            </div>

            <p class="text-xs text-foreground/50">{{ $t('map.editor.coordinateHint') }}</p>
          </section>
        </div>

        <div class="flex flex-col gap-6">
          <section class="flex flex-col gap-3 rounded-xl border border-accent/20 bg-background p-4">
            <h2 class="text-lg font-bold">{{ $t('map.editor.preview') }}</h2>

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
                  <span class="text-base">{{ form.title || '' }}</span>
                </div>
                <p v-if="form.subtitle" class="text-xs font-light italic text-justify">
                  {{ form.subtitle }}
                </p>
              </div>
            </div>
          </section>

          <section class="flex flex-col gap-3 rounded-xl border border-accent/20 bg-background p-4">
            <h2 class="text-lg font-bold">{{ $t('map.editor.coordinates') }}</h2>

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
          </section>
        </div>
      </div>
    </div>

    <div v-else class="flex flex-col items-start gap-3">
      <p class="text-red-400">{{ $t('editor.notFound') }}</p>
      <router-link to="/editor/travel" class="rounded-full border border-dotted px-3 py-1">
        {{ $t('editor.backToList') }}
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

import LocationMapPicker from '@/components/map/LocationMapPicker.vue'
import { getMapLocationBySlug, upsertMapLocation } from '@/services/travelersMap'
import {
  formatZoom,
  isValidLatitude,
  isValidLongitude,
  parseZoom,
  zoomToResolution
} from '@/helpers/map'
import { COLOR_KEYS, TEXT_COLOR_KEYS, parseTailwindColor } from '@/helpers/mapPalette'
import { useToast } from '@/composables/useToast'
import { useUnsavedGuard } from '@/composables/useUnsavedGuard'

import type { MapLocation } from '@/interfaces/TravelersMap'

const route = useRoute()
const { t } = useI18n()
const toast = useToast()

const location = ref<MapLocation | null>(null)
const loading = ref(true)
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

useUnsavedGuard(dirty)

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
  location.value = null

  try {
    location.value = await getMapLocationBySlug(slug.value)

    if (location.value) {
      applyLocation(location.value)
    }
  } catch {
    toast.error(t('map.editor.failedLoad'))
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

function onKeyDown(event: KeyboardEvent): void {
  if (!(event.ctrlKey || event.metaKey) || event.key.toLowerCase() !== 's') return
  if (event.repeat) return

  event.preventDefault()
  void save()
}

watch(slug, load, { immediate: true })

onMounted(() => window.addEventListener('keydown', onKeyDown))
onUnmounted(() => window.removeEventListener('keydown', onKeyDown))
</script>
