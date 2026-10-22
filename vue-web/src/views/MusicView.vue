<template>
  <div class="flex flex-col gap-4">
    <div>
      <p class="text-2xl font-bold">{{ $t('music.title') }}<CursorBlink /></p>
    </div>

    <div class="flex flex-col md:flex-row gap-4">
      <div class="flex flex-col gap-2 pb-4 w-full md:w-2/3 mx-auto">
        <MusicPlayer
          :v="getCurrentSong(current)?.v || ''"
          :title="getCurrentSong(current)?.title"
          :artist="getCurrentSong(current)?.artist"
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
              :disabled="playlist && current >= playlist.length - 1"
              @click="changeSong(1)"
            >
              <Icon class="mx-auto" icon="mynaui:skip-forward" height="1.5rem"></Icon>
            </button>
          </template>
        </MusicPlayer>

        <div class="pt-4 mx-auto">
          <p class="text-lg text-center font-bold">{{ $t('music.meaning') }}</p>
          <p class="text-justify">
            {{ getCurrentSong(current)?.description || $t('music.nothingSpecified') }}<CursorBlink />
          </p>
        </div>
      </div>

      <div class="flex flex-col gap-2 w-full md:w-1/3 mr-auto">
        <div>
          <p class="text-lg font-bold">{{ $t('music.playlist') }}</p>
        </div>

        <div class="max-h-[50vh] overflow-auto scrollbar-thumb-accent scrollbar-track-background scrollbar-thin scroll-pb-6">
          <div
            class="flex border-b py-2 rounded-full max-w-full overflow-x-scroll cursor-pointer scrollbar-0"
            :class="current === idx ? ['bg-accent'] : []"
            v-for="(song, idx) in playlist"
            :key="idx"
            @click="changeSong(idx, true)"
          >
            <p class="px-4 text-nowrap">{{ idx + 1 }}. {{ song.title }} - {{ song.artist }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { getMusic } from '@/services/music'

import type { Music } from '@/interfaces/Music'

const current = ref(0)

const response = getMusic()
const playlist = response.data

function getCurrentSong(pos: number): Music | void {
  if (!playlist?.value) return

  if (pos >= 0 && pos < playlist.value.length) return playlist.value[pos]
}

function changeSong(to: number, relative: boolean = false): void {
  if (relative) {
    current.value = to
  } else {
    current.value += to
  }
}
</script>
