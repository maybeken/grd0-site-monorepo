<template>
  <div class="rounded-xl bg-shade p-4">
    <MusicPlayer :v="$props.block.v" :title="$props.block.title" :artist="$props.block.artist">
    </MusicPlayer>
  </div>
</template>

<script setup lang="ts">
import type { ArticleBlockMusic } from '@/interfaces/Article'

interface Props {
  block: ArticleBlockMusic
}

const $props = defineProps<Props>()
</script>
