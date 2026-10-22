<template>
  <div class="flex flex-col gap-6">
    <template v-for="unit in displayUnits" :key="unitKey(unit)">
      <ImageGrid v-if="unit.kind === 'grid'" :blocks="unit.blocks"></ImageGrid>
      <div v-else :data-block-id="unit.block.id">
        <MarkdownBlock
          v-if="unit.block.type === 'markdown'"
          :block="unit.block"
          :eager="$props.eager"
        ></MarkdownBlock>
        <ImageBlock v-else-if="unit.block.type === 'image'" :block="unit.block"></ImageBlock>
        <MusicBlock v-else-if="unit.block.type === 'music'" :block="unit.block"></MusicBlock>
        <ExcalidrawBlock
          v-else-if="unit.block.type === 'excalidraw'"
          :key="`${unit.block.id}-${$props.version ?? ''}`"
          :block="unit.block"
          :uri="$props.uri"
          :draft="$props.draft"
        ></ExcalidrawBlock>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import ImageGrid from './blocks/ImageGrid.vue'

import type { ArticleBlock, ArticleBlockImage } from '@/interfaces/Article'

interface Props {
  blocks: ArticleBlock[]
  uri: string
  draft?: boolean
  version?: string | number
  eager?: boolean
}

type DisplayUnit =
  | { kind: 'single'; block: ArticleBlock }
  | { kind: 'grid'; blocks: ArticleBlockImage[] }

const $props = defineProps<Props>()

const displayUnits = computed<DisplayUnit[]>(() => {
  const units: DisplayUnit[] = []
  let run: ArticleBlockImage[] = []

  const flush = () => {
    if (run.length === 1) {
      units.push({ kind: 'single', block: run[0] })
    } else if (run.length > 1) {
      units.push({ kind: 'grid', blocks: run })
    }

    run = []
  }

  for (const block of $props.blocks) {
    if (block.type === 'image') {
      run.push(block)
      continue
    }

    flush()
    units.push({ kind: 'single', block })
  }

  flush()

  return units
})

const unitKey = (unit: DisplayUnit) =>
  unit.kind === 'grid' ? unit.blocks.map((block) => block.id).join('-') : unit.block.id
</script>
