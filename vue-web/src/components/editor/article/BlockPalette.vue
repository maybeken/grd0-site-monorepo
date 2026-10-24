<template>
  <div class="flex flex-wrap gap-2">
    <EditorButton
      v-for="item in items"
      :key="item.type"
      variant="secondary"
      size="sm"
      @click="$emit('add', item.type)"
    >
      {{ $t(item.label) }}
    </EditorButton>
  </div>
</template>

<script setup lang="ts">
import EditorButton from '@/components/editor/ui/EditorButton.vue'

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
