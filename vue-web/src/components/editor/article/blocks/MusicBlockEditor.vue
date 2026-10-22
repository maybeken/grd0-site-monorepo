<template>
  <div class="flex flex-col gap-3">
    <div v-if="$props.block.v" class="flex flex-col gap-2 rounded-xl bg-shade p-3">
      <MusicPlayer
        variant="compact"
        :v="$props.block.v"
        :title="$props.block.title"
        :artist="$props.block.artist"
      ></MusicPlayer>
      <button class="rounded-full border border-dotted px-3 py-1" @click="picker = !picker">
        {{ $t('article.editor.changeMusic') }}
      </button>
    </div>

    <div v-else class="flex flex-wrap items-center gap-2">
      <p class="text-secondary">{{ $t('article.editor.noSelection') }}</p>
      <button class="rounded-full border border-dotted px-3 py-1" @click="picker = true">
        {{ $t('article.editor.chooseMusic') }}
      </button>
    </div>

    <MusicPicker v-if="picker" @select="onSelect" @close="picker = false"></MusicPicker>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

import MusicPlayer from '@/components/music/MusicPlayer.vue'

import type { ArticleBlockMusic } from '@/interfaces/Article'

interface Props {
  block: ArticleBlockMusic
}

const $props = defineProps<Props>()

const picker = ref(false)

function onSelect(song: { v: string; title: string; artist?: string }): void {
  $props.block.v = song.v
  $props.block.title = song.title
  $props.block.artist = song.artist
  picker.value = false
}
</script>
