<template>
  <div class="flex flex-col gap-6">
    <header class="flex flex-col gap-2">
      <h1 class="text-2xl font-bold">{{ $t('map.editor.locationsTitle') }}<CursorBlink /></h1>
      <p class="text-secondary">{{ $t('map.editor.locationsHint') }}</p>
    </header>

    <div v-if="!show_form" class="flex justify-end">
      <button
        class="rounded bg-accent px-3 py-1 text-sm font-semibold text-background hover:brightness-110"
        @click="openForm"
      >
        {{ $t('map.editor.newLocation') }}
      </button>
    </div>

    <form
      v-else
      class="flex flex-col gap-4 rounded-xl border border-accent/20 bg-background p-4"
      @submit.prevent="create"
    >
      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold" for="new-map-location-title">
          {{ $t('map.editor.titleField') }}
        </label>
        <input
          id="new-map-location-title"
          ref="title_input"
          v-model="new_title"
          class="rounded border border-accent/40 bg-background px-2 py-1 text-sm"
          :placeholder="$t('map.editor.titlePlaceholder')"
          :disabled="creating"
        />
        <p v-if="generated_slug" class="font-mono text-xs text-foreground/40">
          {{ $t('map.editor.slugField') }}: {{ generated_slug }}
        </p>
      </div>

      <p v-if="collision" class="text-sm text-red-400">
        {{ $t('map.editor.slugCollision', { slug: generated_slug }) }}
      </p>

      <div class="flex justify-end gap-2">
        <button
          type="button"
          class="rounded border border-accent/40 px-3 py-1 text-sm hover:bg-accent/10"
          :disabled="creating"
          @click="closeForm"
        >
          {{ $t('common.cancel') }}
        </button>
        <button
          type="submit"
          class="rounded bg-accent px-3 py-1 text-sm font-semibold text-background hover:brightness-110 disabled:opacity-50"
          :disabled="creating || !generated_slug"
        >
          {{ creating ? $t('common.saving') : $t('map.editor.create') }}
        </button>
      </div>
    </form>

    <div v-if="!show_form" class="flex flex-col gap-2">
      <EditorEntityRow
        v-for="location in locations"
        :key="location.slug"
        :title="location.title"
        :to="`/editor/travel/${location.slug}`"
        @delete="remove(location)"
      >
        <span v-if="location.subtitle" class="truncate text-xs text-foreground/60">
          {{ location.subtitle }}
        </span>
        <span class="flex items-center gap-2">
          <span class="truncate font-mono text-xs text-foreground/40">{{ location.slug }}</span>
          <span
            class="rounded-full px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wide"
            :class="
              location.published
                ? 'bg-emerald-600/20 text-emerald-400'
                : 'bg-orange-500/20 text-orange-400'
            "
          >
            {{ location.published ? $t('map.editor.published') : $t('map.editor.draft') }}
          </span>
        </span>
      </EditorEntityRow>

      <p v-if="loading" class="py-4 text-center text-sm text-foreground/40">
        {{ $t('common.loading') }}
      </p>
      <p
        v-else-if="!locations.length"
        class="py-4 text-center text-sm text-foreground/40"
      >
        {{ $t('map.editor.noLocations') }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import EditorEntityRow from '@/components/editor/EditorEntityRow.vue'
import {
  deleteMapLocation,
  getAllMapLocationsRaw,
  toKebabCase,
  upsertMapLocation
} from '@/services/travelersMap'
import { useToast } from '@/composables/useToast'

import type { MapLocation } from '@/interfaces/TravelersMap'

const { t } = useI18n()
const toast = useToast()
const router = useRouter()

const locations = ref<MapLocation[]>([])
const loading = ref(true)
const show_form = ref(false)
const creating = ref(false)
const new_title = ref('')
const title_input = ref<HTMLInputElement | null>(null)

const generated_slug = computed(() => toKebabCase(new_title.value))
const collision = computed(
  () =>
    !!generated_slug.value &&
    locations.value.some((location) => location.slug === generated_slug.value)
)

async function fetchLocations(): Promise<void> {
  loading.value = true

  try {
    const response = await getAllMapLocationsRaw()

    locations.value = Array.isArray(response) ? response : []
  } catch {
    toast.error(t('map.editor.failedLoad'))
  } finally {
    loading.value = false
  }
}

async function openForm(): Promise<void> {
  new_title.value = ''
  show_form.value = true

  await nextTick()
  title_input.value?.focus()
}

function closeForm(): void {
  show_form.value = false
  new_title.value = ''
}

async function create(): Promise<void> {
  const slug = generated_slug.value
  const title = new_title.value.trim()

  if (!slug || !title || creating.value) return

  if (locations.value.some((location) => location.slug === slug)) {
    toast.error(t('map.editor.slugCollision', { slug }))
    return
  }

  creating.value = true

  try {
    await upsertMapLocation({
      slug,
      title,
      icon: '',
      color: 'default',
      pos: [0, 0],
      published: false
    })

    show_form.value = false
    new_title.value = ''
    await router.push(`/editor/travel/${slug}`)
  } catch {
    toast.error(t('map.editor.failedCreate'))
  } finally {
    creating.value = false
  }
}

async function remove(location: MapLocation): Promise<void> {
  if (!location.slug || !confirm(t('map.editor.confirmDelete', { title: location.title }))) {
    return
  }

  try {
    await deleteMapLocation(location.slug)
    locations.value = locations.value.filter((item) => item.slug !== location.slug)
  } catch {
    toast.error(t('map.editor.failedDelete'))
  }
}

onMounted(fetchLocations)
</script>
