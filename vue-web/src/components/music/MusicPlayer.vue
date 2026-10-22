<template>
  <div class="flex flex-col" :class="$props.variant === 'compact' ? 'gap-3' : 'gap-2'">
    <div class="hidden">
      <YouTubeEmbed
        :v="$props.v"
        :play="play"
        :seek="seek"
        @onReady="onReady"
        @initialDelivery="onMessage"
        @infoDelivery="onMessage"
      ></YouTubeEmbed>
    </div>

    <template v-if="$props.variant === 'compact'">
      <div class="flex items-center gap-3 sm:gap-4">
        <img
          v-if="$props.v"
          class="h-16 w-16 shrink-0 rounded-lg object-cover sm:h-22 sm:w-22"
          :src="`https://i.ytimg.com/vi_webp/${$props.v}/maxresdefault.webp`"
          loading="lazy"
        />

        <div class="flex min-w-0 flex-1 flex-col">
          <p class="truncate text-base font-black sm:text-lg">{{ displayTitle }}</p>
          <p class="truncate text-sm text-secondary">{{ displayArtist }}</p>
        </div>

        <div class="hidden shrink-0 text-sm text-secondary tabular-nums sm:block">
          {{ formatCompactTime(current_time) }} / {{ formatCompactTime(duration) }}
        </div>

        <button
          class="h-10 w-10 shrink-0 rounded-full px-2 border border-white border-dotted bg-black z-10 disabled:brightness-50 sm:h-12 sm:w-12"
          :disabled="loading || !$props.v"
          @click="togglePlay"
        >
          <Icon class="mx-auto" v-if="!play" icon="mynaui:play" height="1.5rem"></Icon>
          <Icon class="mx-auto" v-else icon="mynaui:pause" height="1.5rem"></Icon>
        </button>
      </div>

      <MusicProgressBar
        :progress="current_time"
        :duration="duration"
        @seek="seekTo"
      ></MusicProgressBar>
    </template>

    <template v-else>
      <div v-if="$props.v" class="flex place-content-center">
        <img
          class="rounded-xl w-full md:w-1/2 lg:w-1/3 aspect-square object-cover"
          :src="`https://i.ytimg.com/vi_webp/${$props.v}/maxresdefault.webp`"
          loading="lazy"
        />
      </div>

      <div class="flex flex-col text-center place-content-center">
        <p class="text-xl font-black">{{ displayTitle }}</p>
        <p class="text-lg">{{ displayArtist }}</p>
      </div>

      <div class="flex place-content-center">
        <div class="w-full md:w-1/2 lg:w-1/3">
          <MusicProgressBar
            :progress="current_time"
            :duration="duration"
            @seek="seekTo"
          ></MusicProgressBar>
        </div>
      </div>

      <div class="flex place-content-center">
        <div>{{ convertToTimeFormat(current_time) }} / {{ convertToTimeFormat(duration) }}</div>
      </div>

      <div class="flex place-content-center">
        <slot name="previous"></slot>
        <button
          class="h-18 w-18 rounded-full px-2 border border-white border-dotted bg-black z-10 disabled:brightness-50"
          :disabled="loading || !$props.v"
          @click="togglePlay"
        >
          <div class="w-full">
            <Icon class="mx-auto" v-if="!play" icon="mynaui:play" height="2.5rem"></Icon>
            <Icon class="mx-auto" v-else icon="mynaui:pause" height="2.5rem"></Icon>
          </div>
        </button>
        <slot name="next"></slot>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'

import type { InitialDeliveryMessage, InfoDeliveryMessage } from '@/interfaces/YouTube'
import { t } from '@/i18n'
import MusicProgressBar from '@/components/music/MusicProgressBar.vue'

interface Props {
  v: string
  title?: string
  artist?: string
  variant?: 'full' | 'compact'
}

const $props = withDefaults(defineProps<Props>(), { variant: 'full' })
const $emit = defineEmits(['ended'])

const play = ref(false)
const loading = ref(true)
const initialized = ref(false)
const current_time = ref(0)
const duration = ref(0)
const seek = ref(-1)
const video_title = ref('')
const video_artist = ref('')

const displayTitle = computed(
  () => $props.title || video_title.value || t('music.unknown')
)
const displayArtist = computed(
  () => $props.artist || video_artist.value || t('music.unknown')
)

function onReady(): void {
  current_time.value = 0
  loading.value = false

  if (initialized.value) play.value = true
}

function onMessage(message: InitialDeliveryMessage | InfoDeliveryMessage): void {
  if ('videoData' in message && message.videoData?.title) {
    video_title.value = message.videoData.title
    video_artist.value = message.videoData.author
  }

  if ('currentTime' in message && typeof message.currentTime === 'number') {
    current_time.value = message.currentTime
  }

  if ('duration' in message && typeof message.duration === 'number') {
    duration.value = message.duration
  }

  if ('playerState' in message && message.playerState === 0 && initialized.value) {
    $emit('ended')
  }
}

function convertToTimeFormat(seconds: number): string {
  const rounded = Math.round(seconds || 0)

  const second = rounded % 60
  const minute = (rounded - second) / 60
  const hour = (rounded - second - minute * 60) / 60 / 60

  const second_str = `${second}`.padStart(2, '0')
  const minute_str = `${minute}`.padStart(2, '0')
  const hour_str = `${hour}`.padStart(2, '0')

  return `${hour_str}:${minute_str}:${second_str}`
}

function formatCompactTime(seconds: number): string {
  const rounded = Math.round(seconds || 0)

  const second = rounded % 60
  const minute = Math.floor(rounded / 60) % 60
  const hour = Math.floor(rounded / 3600)

  const second_str = `${second}`.padStart(2, '0')
  const minute_str = `${minute}`.padStart(2, '0')

  return hour > 0 ? `${hour}:${minute_str}:${second_str}` : `${minute}:${second_str}`
}

function seekTo(percentage: number): void {
  seek.value = duration.value * percentage
}

function togglePlay(): void {
  initialized.value = true
  play.value = !play.value
}

watch(
  () => $props.v,
  () => {
    play.value = false
    loading.value = true
    current_time.value = 0
    duration.value = 0
    seek.value = -1
    video_title.value = ''
    video_artist.value = ''
  }
)
</script>
