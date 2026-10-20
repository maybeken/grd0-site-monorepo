<template>
  <div class="flex flex-wrap gap-2">
    <button
      v-for="item in items"
      :key="item.type"
      class="rounded-full border border-dotted px-4 py-1"
      @click="$emit('add', item.type)"
    >
      {{ $t(item.label) }}
    </button>
  </div>
</template>

<script setup lang="ts">
import type { ArticleBlockType } from '@/interfaces/Article'

const $emit = defineEmits<{
  add: [type: ArticleBlockType]
}>()

const items: { type: ArticleBlockType; label: string }[] = [
  { type: 'markdown', label: 'article.editor.blockMarkdown' },
  { type: 'image', label: 'article.editor.blockImage' },
  { type: 'music', label: 'article.editor.blockMusic' },
  { type: 'excalidraw', label: 'article.editor.blockDrawing' }
]
</script>
