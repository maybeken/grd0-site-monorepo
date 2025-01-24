<template>
  <div>
    <p class="text-2xl">Music Matters</p>

    <div class="hidden">
      <YouTube :v="getCurrentSong(current)" :play="play" @onReady="onYouTubeReady" @initialDelivery="youTubeMessageHandler" @infoDelivery="youTubeMessageHandler"></YouTube>
    </div>

    <div class="flex flex-col gap-2 w-full mx-auto">
      <div class="flex place-content-center">
        <img class="rounded-xl w-full md:w-1/2 lg:w-1/3 aspect-square object-cover" :src="`https://i.ytimg.com/vi_webp/${getCurrentSong(current)}/sddefault.webp`" loading="lazy" />
      </div>
      <div class="flex flex-col text-center place-content-center">
        <p class="text-xl font-black">{{ title ?? playlist[current].title }}</p>
        <p class="text-lg">{{ playlist[current].artist ?? "Unknown" }}</p>
      </div>
      <div class="flex place-content-center">
        <div class="relative w-full md:w-1/2 lg:w-1/3">
          <div
            class="absolute left-0 z-20 h-2 rounded-2xl bg-gradient-to-r from-indigo-500 from-0% via-sky-500 via-50% to-emerald-500 to-100%"
            :style="`width: ${current_time/duration*100}%`"
          ></div>
          <div class="absolute left-0 z-10 h-2 rounded-2xl bg-accent w-full"></div>
        </div>
      </div>
      <div class="flex place-content-center"> 
        <div>{{ convertToTimeFormat(current_time) }} / {{ convertToTimeFormat(duration) }}</div>
      </div>
      <div class="flex place-content-center">
        <button
          class="my-auto h-12 w-12 rounded-full px-2 border border-white border-dotted disabled:opacity-50"
          :disabled="current == 0"
          @click="changeSong(-1)"
        >
          <Icon class="mx-auto" icon="mynaui:skip-back" height="1.5rem"></Icon>
        </button>
        <button
          class="h-16 w-16 rounded-full px-2 border border-white border-dotted disabled:opacity-50"
          :disabled="loading"
          @click="play = !play"
        >
          <div class="w-full">
            <Icon class="mx-auto" v-if="!play" icon="mynaui:play" height="2.5rem"></Icon>
            <Icon class="mx-auto" v-else icon="mynaui:pause" height="2.5rem"></Icon>
          </div>
        </button>
        <button
          class="my-auto h-12 w-12 rounded-full px-2 border border-white border-dotted disabled:opacity-50"
          :disabled="current >= playlist.length - 1"
          @click="changeSong(1)"
        >
          <Icon class="mx-auto" icon="mynaui:skip-forward" height="1.5rem"></Icon>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import type { InitialDeliveryMessage } from '@/interfaces/YouTube';

const loading = ref(true);
const play = ref(false);
const title = ref('');
const author = ref('');
const current_time = ref(0);
const duration = ref(0);
const current = ref(0);

const playlist = [
  { title: 'Tubthumping', artist: 'Chumbawamba', v: 'yW5oTzftgjY' },
  { title: 'Virtual Insanity', artist: 'Jamiroquai', v: 'GieQq3eWSnE' },
];

function getCurrentSong(pos: number): string | void {
  if (pos >= 0 && pos < playlist.length) return playlist[pos].v;
}

function onYouTubeReady(): void {
  current_time.value = 0;
  loading.value = false;
}

function youTubeMessageHandler(message: InitialDeliveryMessage): void {
  if (message.videoData?.title) {
    title.value = message.videoData.title;
    author.value = message.videoData.author;
  }

  if (message.currentTime) {
    current_time.value = message.currentTime;
  }

  if (message.duration) {
    duration.value = message.duration;
  }
}

function convertToTimeFormat(seconds: number): string {
  const rounded = Math.round(seconds);

  const second = rounded % 60;
  const minute = (rounded - second) / 60;
  const hour = (rounded - second - minute * 60) / 60 / 60;

  const second_str = `${second}`.padStart(2, '0');
  const minute_str = `${minute}`.padStart(2, '0');
  const hour_str = `${hour}`.padStart(2, '0');

  return `${hour_str}:${minute_str}:${second_str}`;
}

function changeSong(to: number): void {
  current.value += to;
  play.value = false;
  loading.value = true;
}
</script>