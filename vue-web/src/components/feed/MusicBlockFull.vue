<template>
  <div class="flex flex-col gap-2 py-2">
    <div class="hidden">
      <YouTubeEmbed
        v-if="block.music?.v"
        :v="block.music.v"
        :play="isPlaying"
        :seek="seek"
        :stop="false"
        @onReady="onYouTubeReady"
        @initialDelivery="youTubeMessageHandler"
        @infoDelivery="youTubeMessageHandler"
      ></YouTubeEmbed>
    </div>

    <div class="flex items-center gap-4">
      <img
        class="rounded-lg w-20 h-20 object-cover"
        :src="`https://i.ytimg.com/vi_webp/${block.music?.v}/maxresdefault.webp`"
        loading="lazy"
      />
      <div class="flex flex-col flex-1">
        <p class="font-bold">{{ block.music?.title }}</p>
        <p class="text-secondary text-sm">{{ block.music?.artist ?? $t('music.unknown') }}</p>
        <div class="mt-1">
          <MusicProgressBar
            :progress="current_time"
            :duration="duration"
            @seek="(pct: number) => (seek = duration * pct)"
          ></MusicProgressBar>
        </div>
        <div class="text-xs text-secondary">
          {{ convertToTimeFormat(current_time) }} / {{ convertToTimeFormat(duration) }}
        </div>
      </div>
      <button
        class="h-12 w-12 rounded-full border border-white border-dotted flex items-center justify-center disabled:brightness-50"
        :disabled="loading"
        @click="togglePlay"
      >
        <Icon v-if="!isPlaying" icon="mynaui:play" height="1.5rem"></Icon>
        <Icon v-else icon="mynaui:pause" height="1.5rem"></Icon>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { usePlayerCoordinator } from '@/stores/playerCoordinator'

import type { FeedBlock } from '@/interfaces/Feed'
import type { InitialDeliveryMessage, InfoDeliveryMessage } from '@/interfaces/YouTube'

const $props = defineProps<{
  block: FeedBlock
}>()

const coordinator = usePlayerCoordinator()

const loading = ref(true)
const initialized = ref(false)
const current_time = ref(0)
const duration = ref(0)
const seek = ref(-1)

const playerId = computed(() => $props.block.id)
const isPlaying = computed(() => coordinator.activePlayerId === playerId.value)

watch(
  () => coordinator.activePlayerId,
  (newActive) => {
    if (newActive !== null && newActive !== playerId.value) {
      initialized.value = false
    }
  }
)

function togglePlay() {
  if (isPlaying.value) {
    coordinator.pause()
  } else {
    initialized.value = true
    coordinator.play(playerId.value)
  }
}

function onYouTubeReady() {
  current_time.value = 0
  loading.value = false
  if (initialized.value) coordinator.play(playerId.value)
}

function youTubeMessageHandler(message: InitialDeliveryMessage | InfoDeliveryMessage) {
  if ('currentTime' in message && message.currentTime) {
    current_time.value = message.currentTime
  }
  if ('duration' in message && message.duration) {
    duration.value = message.duration
  }
  if ('playerState' in message && message.playerState == 0) {
    coordinator.pause()
  }
}

function convertToTimeFormat(seconds: number): string {
  const rounded = Math.round(seconds)
  const second = rounded % 60
  const minute = (rounded - second) / 60
  const hour = (rounded - second - minute * 60) / 60 / 60
  return `${String(hour).padStart(2, '0')}:${String(minute).padStart(2, '0')}:${String(second).padStart(2, '0')}`
}
</script>
