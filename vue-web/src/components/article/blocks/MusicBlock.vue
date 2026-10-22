<template>
  <div class="rounded-xl bg-shade p-3">
    <MusicPlayer
      variant="compact"
      :v="$props.block.v"
      :title="$props.block.title"
      :artist="$props.block.artist"
    >
    </MusicPlayer>
  </div>
</template>

<script setup lang="ts">
import MusicPlayer from '@/components/music/MusicPlayer.vue'

import type { ArticleBlockMusic } from '@/interfaces/Article'

interface Props {
  block: ArticleBlockMusic
}

const $props = defineProps<Props>()
</script>
