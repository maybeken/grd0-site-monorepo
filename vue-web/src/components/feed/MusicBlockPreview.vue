<template>
  <div class="flex items-center gap-4 py-2">
    <img
      class="rounded-lg w-16 h-16 object-cover"
      :src="`https://i.ytimg.com/vi_webp/${block.music?.v}/default.webp`"
      loading="lazy"
    />
    <div>
      <p class="font-bold">{{ block.music?.title }}</p>
      <p class="text-secondary text-sm">{{ block.music?.artist ?? $t('music.unknown') }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { FeedBlock } from '@/interfaces/Feed'

defineProps<{
  block: FeedBlock
}>()
</script>
