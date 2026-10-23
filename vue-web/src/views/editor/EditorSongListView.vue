<template>
  <div class="flex flex-col gap-6">
    <header class="flex flex-col gap-2">
      <h1 class="text-2xl font-bold">{{ $t('editor.songs') }}<CursorBlink /></h1>
      <p class="text-secondary">{{ $t('music.editor.songsHint') }}</p>
    </header>

    <div v-if="!show_form && !editing_v" class="flex justify-end">
      <button
        class="rounded bg-accent px-3 py-1 text-sm font-semibold text-background hover:brightness-110"
        @click="show_form = true"
      >
        {{ $t('music.editor.newSong') }}
      </button>
    </div>

    <MusicEntryForm
      v-if="show_form"
      :catalogue="songs"
      @saved="onSaved"
      @cancel="show_form = false"
      @edit="onEditExisting"
    ></MusicEntryForm>

    <p v-if="loading" class="py-4 text-sm text-foreground/40">{{ $t('common.loading') }}</p>

    <template v-else-if="!show_form">
      <p v-if="load_failed" class="py-4 text-sm text-red-400">
        {{ $t('music.editor.failedLoadSongs') }}
      </p>

      <p v-else-if="membership_failed" class="text-sm text-foreground/50">
        {{ $t('music.editor.failedMembership') }}
      </p>

      <div v-if="!load_failed" class="flex flex-col gap-2">
        <div v-for="song in songs" :key="song.v">
          <MusicEntryForm
            v-if="editing_v === song.v"
            :entry="song"
            :membership="membership(song.v)"
            :deletable="true"
            @cancel="editing_v = ''"
            @saved="onSaved"
            @delete="deleteSong(song)"
          ></MusicEntryForm>

          <div
            v-else
            class="flex items-center gap-3 rounded-lg border border-accent/20 bg-background p-3"
          >
            <img
              class="h-12 w-12 shrink-0 rounded object-cover"
              :src="`https://i.ytimg.com/vi_webp/${song.v}/default.webp`"
              loading="lazy"
            />

            <div class="flex min-w-0 flex-1 flex-col">
              <span class="truncate text-sm font-semibold">{{ song.title || song.v }}</span>
              <span class="truncate text-xs text-foreground/60">
                {{ song.artist || $t('music.unknown') }}
              </span>
              <span v-if="song.description" class="line-clamp-1 text-xs text-foreground/50">
                {{ song.description }}
              </span>
              <span class="text-xs text-foreground/40">
                {{ $t('music.editor.playlistCount', { count: membership(song.v) }) }}
              </span>
            </div>

            <div class="flex shrink-0 gap-2">
              <button
                class="rounded border border-accent/40 px-2 py-1 text-xs hover:bg-accent/10 disabled:opacity-40"
                :disabled="!!editing_v"
                @click="editing_v = song.v"
              >
                {{ $t('common.edit') }}
              </button>
              <button
                class="rounded border border-red-500/40 px-2 py-1 text-xs text-red-400 hover:bg-red-500/10 disabled:opacity-40"
                :disabled="!!editing_v"
                @click="confirmDelete(song)"
              >
                {{ $t('common.delete') }}
              </button>
            </div>
          </div>
        </div>

        <p v-if="!songs.length" class="py-4 text-center text-sm text-foreground/40">
          {{ $t('music.editor.noSongsAvailable') }}
        </p>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import MusicEntryForm from '@/components/music/MusicEntryForm.vue'
import { deleteMusicRaw, getMusicRaw, getPlaylistsRaw } from '@/services/music'
import { useToast } from '@/composables/useToast'

import type { Music } from '@/interfaces/Music'

const { t } = useI18n()
const toast = useToast()

const songs = ref<Music[]>([])
const counts = ref<Record<string, number>>({})
const loading = ref(true)
const load_failed = ref(false)
const membership_failed = ref(false)
const show_form = ref(false)
const editing_v = ref('')

function membership(v: string): number {
  return counts.value[v] ?? 0
}

async function load(initial = false): Promise<void> {
  if (initial) loading.value = true

  const [music_result, playlist_result] = await Promise.allSettled([
    getMusicRaw(),
    getPlaylistsRaw()
  ])

  if (music_result.status === 'fulfilled' && Array.isArray(music_result.value)) {
    songs.value = music_result.value
    load_failed.value = false
  } else {
    load_failed.value = true
    toast.error(t('music.editor.failedLoadSongs'))
  }

  if (playlist_result.status === 'fulfilled' && Array.isArray(playlist_result.value)) {
    const next: Record<string, number> = {}

    for (const playlist of playlist_result.value) {
      for (const song of playlist.songs ?? []) {
        if (!song.v) continue
        next[song.v] = (next[song.v] ?? 0) + 1
      }
    }

    counts.value = next
    membership_failed.value = false
  } else {
    counts.value = {}
    membership_failed.value = true
    toast.error(t('music.editor.failedMembership'))
  }

  loading.value = false
}

function onEditExisting(v: string): void {
  show_form.value = false
  editing_v.value = v
}

async function onSaved(): Promise<void> {
  show_form.value = false
  editing_v.value = ''
  await load()
}

async function deleteSong(song: Music): Promise<void> {
  try {
    await deleteMusicRaw(song.v)

    if (editing_v.value === song.v) editing_v.value = ''
    await load()
  } catch {
    toast.error(t('music.editor.failedDeleteSong'))
  }
}

function confirmDelete(song: Music): void {
  if (
    !confirm(
      t('music.editor.confirmDeleteSong', {
        title: song.title || song.v,
        count: membership(song.v)
      })
    )
  ) {
    return
  }

  void deleteSong(song)
}

onMounted(() => load(true))
</script>
