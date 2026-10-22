<template>
  <div class="image-grid">
    <ImageBlock
      v-for="block in $props.blocks"
      :key="block.id"
      :block="block"
      :data-block-id="block.id"
    ></ImageBlock>
  </div>
</template>

<script setup lang="ts">
import ImageBlock from './ImageBlock.vue'

import type { ArticleBlockImage } from '@/interfaces/Article'

interface Props {
  blocks: ArticleBlockImage[]
}

const $props = defineProps<Props>()
</script>
