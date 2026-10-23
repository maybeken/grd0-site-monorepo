<template>
  <div class="flex flex-col gap-4">
    <div>
      <p class="text-2xl font-bold">{{ $t('music.title') }}<CursorBlink /></p>
    </div>

    <div class="flex flex-col md:flex-row gap-4">
      <div class="flex flex-col gap-2 pb-4 w-full md:w-2/3 mx-auto">
        <MusicPlayer
          :v="current_song?.v || ''"
          :title="current_song?.title"
          :artist="current_song?.artist"
          @ended="changeSong(1)"
        >
          <template #previous>
            <button
              class="my-auto h-12 w-16 rounded-l-full -mr-2 px-2 border border-white border-dotted disabled:brightness-50"
              :disabled="current == 0"
              @click="changeSong(-1)"
            >
              <Icon class="mx-auto" icon="mynaui:skip-back" height="1.5rem"></Icon>
            </button>
          </template>
          <template #next>
            <button
              class="my-auto h-12 w-16 rounded-r-full -ml-2 px-2 border border-white border-dotted disabled:brightness-50"
              :disabled="current >= songs.length - 1"
              @click="changeSong(1)"
            >
              <Icon class="mx-auto" icon="mynaui:skip-forward" height="1.5rem"></Icon>
            </button>
          </template>
        </MusicPlayer>

        <div class="pt-4 mx-auto">
          <p class="text-lg text-center font-bold">{{ $t('music.meaning') }}</p>
          <p class="text-justify">
            {{ current_song?.description || $t('music.nothingSpecified') }}<CursorBlink />
          </p>
        </div>
      </div>

      <div class="flex flex-col gap-2 w-full md:w-1/3 mr-auto">
        <PlaylistSelector
          v-if="playlists && playlists.length > 0"
          :playlists="playlists"
          :selected="selected_playlist"
          @select="selectPlaylist"
        ></PlaylistSelector>

        <div>
          <p class="text-lg font-bold">{{ $t('music.playlist') }}</p>
        </div>

        <div
          class="max-h-[50vh] overflow-auto scrollbar-thumb-accent scrollbar-track-background scrollbar-thin scroll-pb-6"
        >
          <div
            class="flex border-b py-2 rounded-full max-w-full overflow-x-scroll cursor-pointer scrollbar-0"
            :class="current === idx ? ['bg-accent'] : []"
            v-for="(song, idx) in songs"
            :key="song.v"
            @click="changeSong(idx, true)"
          >
            <p class="px-4 text-nowrap">{{ idx + 1 }}. {{ song.title }} - {{ song.artist }}</p>
          </div>
          <p v-if="songs.length === 0" class="px-4 py-2 text-secondary">
            {{ $t('music.emptyPlaylist') }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { getPlaylists } from '@/services/music'

import type { Music } from '@/interfaces/Music'

const current = ref(0)
const selected_playlist = ref('')

const playlist_response = getPlaylists()
const playlists = playlist_response.data

watch(
  playlists,
  (value) => {
    const list = value ?? []

    if (!list.some((playlist) => playlist.id === selected_playlist.value)) {
      selected_playlist.value = list[0]?.id ?? ''
    }
  },
  { immediate: true }
)

const songs = computed<Music[]>(() => {
  const selected = (playlists.value ?? []).find(
    (playlist) => playlist.id === selected_playlist.value
  )

  return selected?.songs ?? []
})

const current_song = computed<Music | undefined>(() => songs.value[current.value])

watch(selected_playlist, () => {
  current.value = 0
})

function selectPlaylist(id: string): void {
  selected_playlist.value = id
}

function changeSong(to: number, relative: boolean = false): void {
  const next = relative ? to : current.value + to

  if (next < 0 || next >= songs.value.length) return

  current.value = next
}
</script>
