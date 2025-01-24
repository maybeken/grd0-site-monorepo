<template>
  <iframe ref="iframe" class="w-full h-fit aspect-video" :src="buildYtUrl($props.v, false, false)" frameborder="0"
    allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
    referrerpolicy="strict-origin-when-cross-origin"></iframe>
</template>

<script setup lang="ts">
import { watch, useTemplateRef, onMounted } from 'vue';

import type { ApiActions, Message } from '@/interfaces/YouTube';

interface Props {
  v: string,
  play: boolean,
  stop: boolean,
};

const $emit = defineEmits(['onReady', 'initialDelivery', 'infoDelivery', 'apiInfoDelivery']);
const $props = defineProps<Props>();

const base_url = 'https://www.youtube-nocookie.com/embed';
const origin = window.location.origin;

const iframe = useTemplateRef('iframe');

function buildYtUrl(vid: string, autoplay: boolean = false, controls: boolean = true) {
  return `${base_url}/${vid}?autoplay=${autoplay ? 1 : 0}&controls=${controls ? 1 : 0}&origin=${origin}&playsinline=1&rel=0&enablejsapi=1`;
}

function sendYouTubeCommand(command: ApiActions): void {
  iframe.value?.contentWindow?.postMessage(JSON.stringify({ event: "command", func: command }), '*');
}

function askYouTubeForMessage(): void {
  if (iframe.value) {
    iframe.value.contentWindow?.postMessage(JSON.stringify({ event: "listening" }), '*');
  }

  setTimeout(askYouTubeForMessage, 100);
}

onMounted(() => {
  window.addEventListener('message', function (msgEvt) {
    const data: Message = JSON.parse(msgEvt.data) || {};

    $emit(data.event, data.info);
  });

  askYouTubeForMessage();
});

watch($props, (newProp) => {
  if (!newProp.stop) {
    sendYouTubeCommand(newProp.play ? 'playVideo' : 'pauseVideo');
  } else {
    sendYouTubeCommand('stopVideo');
  }
})
</script>