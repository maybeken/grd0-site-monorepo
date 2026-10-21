<template>
  <div
    class="flex flex-col gap-2 rounded-xl border border-dotted p-3"
    :class="$props.active ? 'border-accent' : ''"
  >
    <div class="flex items-center justify-between gap-2">
      <span class="text-sm text-secondary">{{ $t(labelKey) }}</span>

      <div class="flex gap-2">
        <button
          class="rounded-full border border-dotted p-1 disabled:opacity-30"
          :title="$t('article.editor.moveUp')"
          :disabled="$props.index <= 0"
          @click="$emit('move-up')"
        >
          <Icon icon="mynaui:arrow-up" height="1.25rem"></Icon>
        </button>
        <button
          class="rounded-full border border-dotted p-1 disabled:opacity-30"
          :title="$t('article.editor.moveDown')"
          :disabled="$props.index >= $props.count - 1"
          @click="$emit('move-down')"
        >
          <Icon icon="mynaui:arrow-down" height="1.25rem"></Icon>
        </button>
        <button
          class="rounded-full border border-dotted p-1"
          :title="$t('article.editor.duplicateBlock')"
          @click="$emit('duplicate')"
        >
          <Icon icon="mynaui:copy" height="1.25rem"></Icon>
        </button>
        <button
          class="rounded-full border border-dotted p-1"
          :title="$t('article.editor.removeBlock')"
          @click="$emit('remove')"
        >
          <Icon icon="mynaui:trash" height="1.25rem"></Icon>
        </button>
      </div>
    </div>

    <slot></slot>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import type { ArticleBlockType } from '@/interfaces/Article'

interface Props {
  type: ArticleBlockType
  index: number
  count: number
  active?: boolean
}

const $props = defineProps<Props>()
const $emit = defineEmits(['move-up', 'move-down', 'duplicate', 'remove'])

const labels: Record<ArticleBlockType, string> = {
  markdown: 'article.editor.blockMarkdown',
  image: 'article.editor.blockImage',
  music: 'article.editor.blockMusic',
  excalidraw: 'article.editor.blockDrawing'
}

const labelKey = computed(() => labels[$props.type])
</script>
