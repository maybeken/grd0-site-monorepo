<template>
  <div class="flex flex-col gap-2 py-2">
    <div class="flex gap-4">
      <div class="h-12 w-12">
        <ProfileIcon :email="block.blog?.author?.email ?? ''"></ProfileIcon>
      </div>
      <div class="my-auto">
        <h3 class="font-bold">{{ block.blog?.title }}</h3>
        <p class="text-secondary text-sm">{{ block.blog?.subtitle }}</p>
      </div>
    </div>
    <div class="max-h-[33vh] overflow-hidden relative">
      <MarkdownDisplay :md="block.blog?.content ?? ''" :lazy_loading="false"></MarkdownDisplay>
      <div
        v-if="!expanded"
        class="absolute bottom-0 left-0 h-16 w-full bg-linear-to-t from-background from-30%"
      ></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { FeedBlock } from '@/interfaces/Feed'

defineProps<{
  block: FeedBlock
  expanded: boolean
}>()
</script>
