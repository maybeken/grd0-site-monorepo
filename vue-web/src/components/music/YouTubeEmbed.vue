<template>
  <iframe
    v-if="$props.v"
    ref="iframe"
    class="w-full h-fit aspect-video"
    :src="buildYtUrl($props.v, false, false)"
    frameborder="0"
    allow="
      accelerometer;
      autoplay;
      clipboard-write;
      encrypted-media;
      gyroscope;
      picture-in-picture;
      web-share;
    "
    referrerpolicy="strict-origin-when-cross-origin"
  ></iframe>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, useTemplateRef, watch } from 'vue'

import type { ApiActions, Message } from '@/interfaces/YouTube'

interface Props {
  v: string
  play?: boolean
  stop?: boolean
  seek?: number
}

const $emit = defineEmits(['onReady', 'initialDelivery', 'infoDelivery', 'apiInfoDelivery'])
const $props = withDefaults(defineProps<Props>(), {
  play: false,
  stop: false,
  seek: -1
})

const base_url = 'https://www.youtube-nocookie.com/embed'
const origin = window.location.origin

const iframe = useTemplateRef('iframe')

let poll_timer: number | undefined

function buildYtUrl(vid: string, autoplay: boolean = false, controls: boolean = true) {
  return `${base_url}/${vid}?autoplay=${autoplay ? 1 : 0}&controls=${controls ? 1 : 0}&origin=${origin}&playsinline=1&rel=0&enablejsapi=1`
}

function sendYouTubeCommand(command: ApiActions, params: any = []): void {
  iframe.value?.contentWindow?.postMessage(
    JSON.stringify({ event: 'command', func: command, args: params }),
    '*'
  )
}

function askYouTubeForMessage(): void {
  iframe.value?.contentWindow?.postMessage(JSON.stringify({ event: 'listening' }), '*')
  poll_timer = window.setTimeout(askYouTubeForMessage, 100)
}

function onMessage(msgEvt: MessageEvent): void {
  if (msgEvt.source !== iframe.value?.contentWindow) return
  if (typeof msgEvt.data !== 'string') return

  let data: Message
  try {
    data = JSON.parse(msgEvt.data) as Message
  } catch {
    return
  }

  if (!data || typeof data.event !== 'string') return

  $emit(data.event, data.info)
}

onMounted(() => {
  window.addEventListener('message', onMessage)
  askYouTubeForMessage()
})

onUnmounted(() => {
  window.removeEventListener('message', onMessage)
  window.clearTimeout(poll_timer)
})

watch(
  () => [$props.v, $props.play, $props.stop, $props.seek],
  () => {
    if ($props.v) {
      if (!$props.stop) {
        sendYouTubeCommand($props.play ? 'playVideo' : 'pauseVideo')
      } else {
        sendYouTubeCommand('stopVideo')
      }

      if ($props.seek > -1) {
        sendYouTubeCommand('seekTo', [Math.round($props.seek), true])
      }
    }
  }
)
</script>
