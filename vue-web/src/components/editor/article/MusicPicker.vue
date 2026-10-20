<template>
  <div class="flex flex-col gap-3 rounded-xl border border-dotted p-3">
    <div class="flex items-center gap-2">
      <span class="font-bold">{{ $t('article.editor.musicTitle') }}</span>
      <button class="ml-auto rounded-full border border-dotted px-3 py-1" @click="$emit('close')">
        {{ $t('common.cancel') }}
      </button>
    </div>

    <p v-if="!loading && (!songs || songs.length === 0)" class="text-sm text-secondary">
      {{ $t('article.editor.noSongs') }}
    </p>

    <div class="flex flex-col max-h-[50vh] overflow-y-auto">
      <button
        v-for="song in songs"
        :key="song.v"
        class="flex items-center gap-3 border-b py-2 text-left hover:bg-shade"
        @click="select(song)"
      >
        <img
          class="h-12 w-12 rounded-lg object-cover"
          :src="`https://i.ytimg.com/vi_webp/${song.v}/default.webp`"
          loading="lazy"
        />
        <span class="flex flex-col">
          <span class="font-bold">{{ song.title }}</span>
          <span class="text-sm text-secondary">{{ song.artist || $t('music.unknown') }}</span>
        </span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getMusic } from '@/services/music'

import type { Music } from '@/interfaces/Music'

const $emit = defineEmits<{
  select: [song: { v: string; title: string; artist?: string }]
  close: []
}>()

const response = getMusic()
const songs = response.data
const loading = response.loading

function select(song: Music): void {
  $emit('select', { v: song.v, title: song.title, artist: song.artist })
}
</script>
