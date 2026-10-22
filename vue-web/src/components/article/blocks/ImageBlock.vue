<template>
  <figure class="flex flex-col gap-2">
    <p v-if="failed || !image_uri" class="rounded-xl border border-dotted p-6 text-center text-secondary">
      {{ $t('article.block.imageUnavailable') }}
    </p>
    <CDNImage
      v-else
      class="rounded-xl w-full object-contain"
      :uri="image_uri"
      :resolution="1280"
      :quality="80"
      @error="failed = true"
    ></CDNImage>

    <figcaption v-if="$props.block.caption" class="article-caption">
      {{ $props.block.caption }}
    </figcaption>
  </figure>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

import type { ArticleBlockImage } from '@/interfaces/Article'

interface Props {
  block: ArticleBlockImage
}

const $props = defineProps<Props>()

const failed = ref(false)

const image_uri = computed(() => {
  if (!$props.block.collection || !$props.block.filename) return ''

  return `${$props.block.collection}/${$props.block.filename}`
})
</script>
