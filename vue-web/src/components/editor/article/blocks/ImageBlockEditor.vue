<template>
  <div class="flex flex-col gap-2">
    <div v-if="preview_uri" class="flex flex-col md:flex-row gap-3">
      <CDNImage
        class="w-full md:w-1/2 rounded-xl"
        :uri="preview_uri"
        :resolution="640"
        :quality="75"
      ></CDNImage>

      <div class="flex flex-col gap-2 w-full">
        <span class="text-sm text-secondary">{{ $props.block.filename }}</span>
        <input
          class="rounded-full border border-dotted px-3 py-1 bg-background"
          type="text"
          :placeholder="$t('article.editor.captionPlaceholder')"
          v-model="$props.block.caption"
        />
        <button class="rounded-full border border-dotted px-3 py-1" @click="picker = !picker">
          {{ $t('article.editor.changeImage') }}
        </button>
      </div>
    </div>

    <div v-else class="flex flex-wrap items-center gap-2">
      <p class="text-secondary">{{ $t('article.editor.noSelection') }}</p>
      <button class="rounded-full border border-dotted px-3 py-1" @click="picker = true">
        {{ $t('article.editor.chooseImage') }}
      </button>
    </div>

    <GalleryPicker
      v-if="picker"
      @select="onSelect"
      @close="picker = false"
    ></GalleryPicker>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

import type { ArticleBlockImage } from '@/interfaces/Article'

interface Props {
  block: ArticleBlockImage
}

const $props = defineProps<Props>()

const picker = ref(false)

const preview_uri = computed(() => {
  if (!$props.block.collection || !$props.block.filename) return ''

  return `${$props.block.collection}/${$props.block.filename}`
})

function onSelect(selection: { collection: string; filename: string }): void {
  $props.block.collection = selection.collection
  $props.block.filename = selection.filename
  picker.value = false
}
</script>
