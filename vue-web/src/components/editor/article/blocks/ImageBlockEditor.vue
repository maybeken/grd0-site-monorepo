<template>
  <div class="flex flex-col gap-2">
    <div v-if="preview_uri" class="flex flex-col gap-3 md:flex-row">
      <CDNImage
        class="w-full rounded-xl md:w-1/2"
        :uri="preview_uri"
        :resolution="640"
        :quality="75"
      ></CDNImage>

      <div class="flex w-full flex-col gap-2">
        <span class="truncate text-sm font-medium text-foreground/60" :title="$props.block.filename">
          {{ $props.block.filename }}
        </span>
        <EditorInput
          v-model="$props.block.caption"
          :placeholder="$t('article.editor.captionPlaceholder')"
        ></EditorInput>
        <EditorButton variant="secondary" size="sm" @click="picker = !picker">
          {{ $t('article.editor.changeImage') }}
        </EditorButton>
      </div>
    </div>

    <div v-else class="flex flex-wrap items-center gap-2">
      <p class="text-lg font-medium text-foreground/60">{{ $t('article.editor.noSelection') }}</p>
      <EditorButton variant="secondary" size="sm" @click="picker = true">
        {{ $t('article.editor.chooseImage') }}
      </EditorButton>
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

import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorInput from '@/components/editor/ui/EditorInput.vue'
import CDNImage from '@/components/shared/CDNImage.vue'
import GalleryPicker from './../GalleryPicker.vue'

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
