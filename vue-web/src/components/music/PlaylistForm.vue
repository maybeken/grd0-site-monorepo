<template>
  <div class="flex flex-col gap-4 bg-background rounded-xl p-4 border border-accent/20">
    <div class="flex items-center justify-between">
      <h3 class="text-lg font-bold">
        {{ is_edit ? $t('music.editor.titleEdit') : $t('music.editor.titleNew') }}
      </h3>
      <div class="flex gap-2">
        <button
          @click="$emit('cancel')"
          class="px-3 py-1 rounded border border-accent/40 text-sm hover:bg-accent/10"
        >
          {{ $t('common.cancel') }}
        </button>
        <button
          @click="save"
          :disabled="saving || !form.title.trim()"
          class="px-3 py-1 rounded bg-accent text-background text-sm font-semibold hover:brightness-110 disabled:opacity-50"
        >
          {{ saving ? $t('common.saving') : $t('common.save') }}
        </button>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">{{ $t('music.editor.titleField') }}</label>
        <input
          v-model="form.title"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
          :placeholder="$t('music.editor.titlePlaceholder')"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">{{ $t('music.editor.sortingField') }}</label>
        <input
          v-model.number="form.sorting"
          type="number"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
        />
      </div>
    </div>

    <div class="flex flex-col gap-1">
      <label class="text-sm font-semibold">{{ $t('music.editor.descriptionField') }}</label>
      <textarea
        v-model="form.description"
        rows="2"
        class="bg-background border border-accent/40 rounded px-2 py-1 text-sm resize-y"
        :placeholder="$t('music.editor.descriptionPlaceholder')"
      />
    </div>

    <div class="flex flex-col gap-2">
      <div class="flex items-center justify-between">
        <label class="text-sm font-semibold">
          {{ $t('music.editor.songsField') }} ({{ selected_ids.length }})
        </label>
        <button
          class="px-2 py-1 rounded text-xs border border-accent/40 hover:bg-accent/10"
          @click="picker = true"
        >
          {{ $t('music.editor.addSong') }}
        </button>
      </div>

      <MusicPicker
        v-if="picker"
        :link_only="true"
        @select="onAddSong"
        @close="picker = false"
      ></MusicPicker>

      <p v-if="!catalogue_songs || catalogue_songs.length === 0" class="text-sm text-secondary">
        {{ $t('music.editor.noSongsAvailable') }}
      </p>

      <div v-else class="flex max-h-72 flex-col overflow-y-auto rounded border border-accent/20">
        <label
          v-for="song in catalogue_songs"
          :key="song.id || song.v"
          class="flex cursor-pointer items-center gap-3 border-b border-accent/10 px-3 py-2 hover:bg-shade"
        >
          <input
            type="checkbox"
            class="accent-accent"
            :checked="song.id ? selected_ids.includes(song.id) : false"
            @change="toggle(song)"
          />
          <img
            class="h-10 w-10 rounded object-cover"
            :src="`https://i.ytimg.com/vi_webp/${song.v}/default.webp`"
            loading="lazy"
          />
          <span class="flex min-w-0 flex-1 flex-col">
            <span class="truncate text-sm font-semibold">{{ song.title }}</span>
            <span class="truncate text-xs text-secondary">
              {{ song.artist || $t('music.unknown') }}
            </span>
          </span>
          <button
            class="shrink-0 rounded border border-red-500/40 px-2 py-1 text-xs text-red-400 hover:bg-red-500/10"
            @click.stop.prevent="removeSong(song)"
          >
            {{ $t('music.editor.removeSong') }}
          </button>
        </label>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import MusicPicker from '@/components/editor/article/MusicPicker.vue'
import { getMusicRaw, deleteMusicRaw, upsertMusicRaw, upsertPlaylistRaw } from '@/services/music'
import { useToast } from '@/composables/useToast'

import type { Music, Playlist } from '@/interfaces/Music'
import { useUnsavedGuard } from '@/composables/useUnsavedGuard'

interface Props {
  playlist?: Playlist | null
}

const props = withDefaults(defineProps<Props>(), {
  playlist: null
})

const $emit = defineEmits<{
  saved: []
  cancel: []
}>()

const { t } = useI18n()
const toast = useToast()

const catalogue_songs = ref<Music[]>([])
const catalogue_loaded = ref(false)
const picker = ref(false)

async function fetchCatalogue(): Promise<void> {
  const data = await getMusicRaw()

  if (Array.isArray(data)) {
    catalogue_songs.value = data
    catalogue_loaded.value = true
  }
}

const catalogue_ready = fetchCatalogue().catch(() => undefined)

const is_edit = computed(() => !!props.playlist?.id)
const saving = ref(false)
const selected_ids = ref<string[]>([])

const form = reactive({
  title: '',
  description: '',
  sorting: 0
})

const baseline = ref('')

function snapshot(): string {
  return JSON.stringify({ form, selected_ids: selected_ids.value })
}

watch(
  () => props.playlist,
  (playlist) => {
    form.title = playlist?.title ?? ''
    form.description = playlist?.description ?? ''
    form.sorting = playlist?.sorting ?? 0
    selected_ids.value = (playlist?.songs ?? [])
      .map((song) => song.id)
      .filter((id): id is string => !!id)

    baseline.value = snapshot()
  },
  { immediate: true }
)

const dirty = computed(() => snapshot() !== baseline.value)

useUnsavedGuard(dirty)

function toggle(song: Music): void {
  if (!song.id) return

  if (selected_ids.value.includes(song.id)) {
    selected_ids.value = selected_ids.value.filter((id) => id !== song.id)
  } else {
    selected_ids.value = [...selected_ids.value, song.id]
  }
}

async function onAddSong(song: { v: string; title: string; artist?: string }): Promise<void> {
  if (!song.v) return

  try {
    await catalogue_ready

    const existing = catalogue_songs.value.find((item) => item.v === song.v)

    if (existing) {
      if (existing.title !== song.title || (existing.artist ?? '') !== (song.artist ?? '')) {
        toast.info(t('music.editor.correctOnSongsPage'))
      }
    } else {
      await upsertMusicRaw({ v: song.v, title: song.title, artist: song.artist })
      await fetchCatalogue()
    }

    const entry = catalogue_songs.value.find((item) => item.v === song.v)

    if (entry?.id && !selected_ids.value.includes(entry.id)) {
      selected_ids.value = [...selected_ids.value, entry.id]
    }

    picker.value = false
  } catch {
    toast.error(t('music.editor.failedAddSong'))
  }
}

async function removeSong(song: Music): Promise<void> {
  if (!song.v || !confirm(t('music.editor.confirmRemoveSong', { title: song.title }))) return

  try {
    await deleteMusicRaw(song.v)

    if (song.id) {
      selected_ids.value = selected_ids.value.filter((id) => id !== song.id)
    }

    await fetchCatalogue()
  } catch {
    toast.error(t('music.editor.failedRemoveSong'))
  }
}

async function save(): Promise<void> {
  if (!form.title.trim() || saving.value) return

  saving.value = true

  try {
    const songs = catalogue_loaded.value
      ? catalogue_songs.value.filter(
          (song) => song.id && selected_ids.value.includes(song.id)
        )
      : (props.playlist?.songs ?? [])

    await upsertPlaylistRaw({
      id: props.playlist?.id,
      title: form.title.trim(),
      description: form.description,
      sorting: form.sorting,
      songs
    })

    baseline.value = snapshot()
    $emit('saved')
  } catch {
    toast.error(t('music.editor.failedSave'))
  } finally {
    saving.value = false
  }
}

function on_save_shortcut(event: KeyboardEvent): void {
  if (!(event.ctrlKey || event.metaKey) || event.key.toLowerCase() !== 's') return

  event.preventDefault()
  void save()
}

onMounted(() => window.addEventListener('keydown', on_save_shortcut))
onBeforeUnmount(() => window.removeEventListener('keydown', on_save_shortcut))
</script>
