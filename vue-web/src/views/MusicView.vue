<template>
  <div class="flex flex-col gap-4">
    <div class="hidden">
      <YouTubeEmbed
        :v="getCurrentSong(current)?.v"
        :play="play"
        :seek="seek"
        @onReady="onYouTubeReady"
        @initialDelivery="youTubeMessageHandler"
        @infoDelivery="youTubeMessageHandler"
      ></YouTubeEmbed>
    </div>

    <div>
      <p class="text-2xl font-bold">{{ $t('music.title') }}<CursorBlink /></p>
    </div>

    <div class="flex flex-col md:flex-row gap-4">
      <div class="flex flex-col gap-2 pb-4 w-full md:w-2/3 mx-auto bg-background">
        <div class="flex place-content-center">
          <img
            class="rounded-xl w-full md:w-1/2 lg:w-1/3 aspect-square object-cover"
            :src="`https://i.ytimg.com/vi_webp/${getCurrentSong(current)?.v}/maxresdefault.webp`"
            loading="lazy"
          />
        </div>

        <div class="flex flex-col text-center place-content-center">
          <p class="text-xl font-black">{{ title ?? getCurrentSong(current)?.title }}</p>
          <p class="text-lg">{{ getCurrentSong(current)?.artist ?? $t('music.unknown') }}</p>
        </div>

        <div class="flex place-content-center">
          <MusicProgressBar
            :progress="current_time"
            :duration="duration"
            @seek="(percentage: number) => (seek = duration * percentage)"
          ></MusicProgressBar>
        </div>

        <div class="flex place-content-center">
          <div>{{ convertToTimeFormat(current_time) }} / {{ convertToTimeFormat(duration) }}</div>
        </div>

        <div class="flex place-content-center">
          <button
            class="my-auto h-12 w-16 rounded-l-full -mr-2 px-2 border border-white border-dotted disabled:brightness-50"
            :disabled="current == 0"
            @click="changeSong(-1)"
          >
            <Icon class="mx-auto" icon="mynaui:skip-back" height="1.5rem"></Icon>
          </button>
          <button
            class="h-18 w-18 rounded-full px-2 border border-white border-dotted bg-black z-10 disabled:brightness-50"
            :disabled="loading"
            @click="playSong"
          >
            <div class="w-full">
              <Icon class="mx-auto" v-if="!play" icon="mynaui:play" height="2.5rem"></Icon>
              <Icon class="mx-auto" v-else icon="mynaui:pause" height="2.5rem"></Icon>
            </div>
          </button>
          <button
            class="my-auto h-12 w-16 rounded-r-full -ml-2 px-2 border border-white border-dotted disabled:brightness-50"
            :disabled="playlist && current >= playlist.length - 1"
            @click="changeSong(1)"
          >
            <Icon class="mx-auto" icon="mynaui:skip-forward" height="1.5rem"></Icon>
          </button>
        </div>

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
import type { InitialDeliveryMessage, InfoDeliveryMessage } from '@/interfaces/YouTube'
import MusicProgressBar from '@/components/music/MusicProgressBar.vue'

const loading = ref(true)
const initialized = ref(false)
const play = ref(false)
const title = ref('')
const author = ref('')
const current_time = ref(0)
const duration = ref(0)
const current = ref(0)
const seek = ref(-1)

const response = getMusic()
const playlist = response.data

function getCurrentSong(pos: number): Music | void {
  if (!playlist?.value) return

  if (pos >= 0 && pos < playlist.value.length) return playlist.value[pos]
}

function onYouTubeReady(): void {
  current_time.value = 0
  loading.value = false

  if (initialized.value) play.value = true
}

function youTubeMessageHandler(message: InitialDeliveryMessage | InfoDeliveryMessage): void {
  if ('videoData' in message && message.videoData?.title) {
    title.value = message.videoData.title
    author.value = message.videoData.author
  }

  if ('currentTime' in message && message.currentTime) {
    current_time.value = message.currentTime
  }

  if ('duration' in message && message.duration) {
    duration.value = message.duration
  }

  if ('playerState' in message) {
    const state = message.playerState

    if (state == 0) {
      changeSong(1)
    }
  }
}

function convertToTimeFormat(seconds: number): string {
  const rounded = Math.round(seconds)

  const second = rounded % 60
  const minute = (rounded - second) / 60
  const hour = (rounded - second - minute * 60) / 60 / 60

  const second_str = `${second}`.padStart(2, '0')
  const minute_str = `${minute}`.padStart(2, '0')
  const hour_str = `${hour}`.padStart(2, '0')

  return `${hour_str}:${minute_str}:${second_str}`
}

function changeSong(to: number, relative: boolean = false): void {
  if (relative) {
    current.value = to
  } else {
    current.value += to
  }

  play.value = false
  loading.value = true
  seek.value = -1
}

function playSong(): void {
  initialized.value = true
  play.value = !play.value
}
</script>
