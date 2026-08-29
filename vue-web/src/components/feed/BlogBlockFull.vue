<template>
  <div class="py-2">
    <h3 v-if="block.blog?.title" class="font-bold text-lg">{{ block.blog?.title }}</h3>
    <hr class="border-accent my-2" />
    <div class="px-2 py-2 text-justify">
      <MarkdownDisplay :md="block.blog?.content ?? ''"></MarkdownDisplay>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { FeedBlock } from '@/interfaces/Feed'

defineProps<{
  block: FeedBlock
}>()
</script>
