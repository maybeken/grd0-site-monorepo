<template>
  <div class="flex flex-col gap-3">
    <div v-if="$props.block.v" class="flex flex-col gap-2 rounded-xl bg-shade p-3">
      <MusicPlayer
        variant="compact"
        :v="$props.block.v"
        :title="$props.block.title"
        :artist="$props.block.artist"
      ></MusicPlayer>
      <EditorButton variant="secondary" size="sm" @click="picker = !picker">
        {{ $t('article.editor.changeMusic') }}
      </EditorButton>
    </div>

    <div v-else class="flex flex-wrap items-center gap-2">
      <p class="text-lg font-medium text-foreground/60">{{ $t('article.editor.noSelection') }}</p>
      <EditorButton variant="secondary" size="sm" @click="picker = true">
        {{ $t('article.editor.chooseMusic') }}
      </EditorButton>
    </div>

    <MusicPicker v-if="picker" @select="onSelect" @close="picker = false"></MusicPicker>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

import EditorButton from '@/components/editor/ui/EditorButton.vue'
import MusicPlayer from '@/components/music/MusicPlayer.vue'
import MusicPicker from './../MusicPicker.vue'

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
