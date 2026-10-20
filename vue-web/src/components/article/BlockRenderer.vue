<template>
  <div class="flex flex-col gap-6">
    <template v-for="block in $props.blocks" :key="block.id">
      <MarkdownBlock v-if="block.type === 'markdown'" :block="block"></MarkdownBlock>
      <ImageBlock v-else-if="block.type === 'image'" :block="block"></ImageBlock>
      <MusicBlock v-else-if="block.type === 'music'" :block="block"></MusicBlock>
      <ExcalidrawBlock
        v-else-if="block.type === 'excalidraw'"
        :key="`${block.id}-${$props.version ?? ''}`"
        :block="block"
        :uri="$props.uri"
        :draft="$props.draft"
      ></ExcalidrawBlock>
    </template>
  </div>
</template>

<script setup lang="ts">
import type { ArticleBlock } from '@/interfaces/Article'

interface Props {
  blocks: ArticleBlock[]
  uri: string
  draft?: boolean
  version?: string | number
}

const $props = defineProps<Props>()
</script>
