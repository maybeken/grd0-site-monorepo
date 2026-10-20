<template>
  <div v-if="loading" class="grid grid-cols-1 lg:grid-cols-2 gap-6 items-start">
    <div class="flex flex-col gap-3">
      <Skeleton h="lg" w="full" :loading="true"></Skeleton>
      <Skeleton h="md" w="2/3" :loading="true"></Skeleton>
      <Skeleton h="xl" w="full" :loading="true"></Skeleton>
    </div>
    <div class="flex flex-col gap-3 rounded-xl border border-dotted p-4">
      <Skeleton h="md" w="1/3" :loading="true"></Skeleton>
      <Skeleton h="xl" w="full" :loading="true"></Skeleton>
      <Skeleton h="lg" w="full" :loading="true"></Skeleton>
    </div>
  </div>

  <div v-else-if="article" class="grid grid-cols-1 lg:grid-cols-2 gap-6 items-start">
    <div class="flex flex-col gap-6 min-w-0">
      <div class="flex flex-col gap-3">
        <input
          class="text-3xl w-full rounded-full px-4 bg-background border border-dotted"
          type="text"
          :placeholder="$t('article.editor.titlePlaceholder')"
          v-model="title"
        />
        <input
          class="text-2xl w-full rounded-full px-4 bg-background border border-dotted"
          type="text"
          :placeholder="$t('article.editor.subtitlePlaceholder')"
          v-model="subtitle"
        />

        <div class="flex flex-wrap items-center gap-2 px-1">
          <label for="article-uri">{{ $t('article.editor.uriLabel') }}</label>
          <span class="text-secondary">/article/</span>
          <input
            id="article-uri"
            class="rounded-full px-3 py-1 bg-background border border-dotted"
            type="text"
            :placeholder="$t('article.editor.uriPlaceholder')"
            v-model="uri"
          />
        </div>

        <div class="flex flex-wrap items-center gap-2 px-1">
          <label for="article-publish">{{ $t('article.editor.publishDate') }}</label>
          <input
            id="article-publish"
            class="rounded-xl px-2 py-1 bg-background border border-dotted"
            type="datetime-local"
            v-model="publish_local"
          />
          <span v-if="!published_at" class="text-secondary">
            {{ $t('article.editor.unpublished') }}
          </span>
        </div>

        <p v-if="conflict" class="text-red-400">{{ $t('article.editor.conflict') }}</p>
        <p v-else-if="error" class="text-red-400">{{ $t(error_message) }}</p>
      </div>

      <div class="flex flex-col gap-3">
        <p class="text-lg font-bold">{{ $t('article.editor.blocks') }}</p>

        <BlockItem
          v-for="(block, index) in blocks"
          :key="block.id"
          :ref="(element) => setBlockElement(block.id, element)"
          :type="block.type"
          :index="index"
          :count="blocks.length"
          @move-up="moveBlock(index, -1)"
          @move-down="moveBlock(index, 1)"
          @remove="removeBlock(index)"
        >
          <MarkdownBlockEditor
            v-if="block.type === 'markdown'"
            :block="block"
          ></MarkdownBlockEditor>
          <ImageBlockEditor v-else-if="block.type === 'image'" :block="block"></ImageBlockEditor>
          <MusicBlockEditor v-else-if="block.type === 'music'" :block="block"></MusicBlockEditor>
          <ExcalidrawBlockEditor
            v-else-if="block.type === 'excalidraw'"
            :ref="(element) => setDrawingRef(block.id, element)"
            :block="block"
            :uri="article?.uri || ''"
            @request-save="saveArticle"
          ></ExcalidrawBlockEditor>
        </BlockItem>
      </div>

      <div class="flex flex-col gap-2">
        <p class="text-lg font-bold">{{ $t('article.editor.addBlock') }}</p>
        <BlockPalette @add="addBlock"></BlockPalette>
      </div>

      <div class="flex flex-wrap gap-2">
        <button
          class="rounded-full border border-dotted px-4 py-1 disabled:opacity-50"
          :disabled="saving"
          aria-keyshortcuts="Control+S Meta+S"
          @click="saveArticle"
        >
          {{ saving ? $t('common.saving') : $t('common.save') }}
        </button>
      </div>
    </div>

    <aside
      class="flex flex-col gap-4 min-w-0 rounded-xl border border-dotted p-4 lg:sticky lg:top-24 lg:max-h-[calc(100vh-7rem)] lg:overflow-y-auto"
    >
      <div class="flex flex-col gap-1">
        <p class="text-lg font-bold">{{ $t('article.editor.preview') }}</p>
        <p class="text-xs text-secondary">{{ $t('article.editor.previewHint') }}</p>
      </div>

      <ArticleHeader :article="preview_article" :draft="true"></ArticleHeader>
      <BlockRenderer
        :blocks="preview_blocks"
        :uri="preview_uri"
        :version="preview_article.updated_at"
        :draft="true"
      ></BlockRenderer>
    </aside>
  </div>

  <div v-else class="flex flex-col gap-2">
    <p v-if="error" class="text-red-400">{{ $t('article.editor.loadFailed') }}</p>
    <p v-else class="text-secondary">{{ $t('article.editor.selectArticle') }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { onBeforeRouteLeave } from 'vue-router'
import dayjs from 'dayjs'

import { createArticleBlock } from '@/helpers/article'
import { useToast } from '@/composables/useToast'
import BlockItem from './BlockItem.vue'
import BlockPalette from './BlockPalette.vue'
import MarkdownBlockEditor from './blocks/MarkdownBlockEditor.vue'
import ImageBlockEditor from './blocks/ImageBlockEditor.vue'
import MusicBlockEditor from './blocks/MusicBlockEditor.vue'
import ExcalidrawBlockEditor from './blocks/ExcalidrawBlockEditor.vue'
import { useArticleDraft } from './useArticleDraft'

import type { Article, ArticleBlock, ArticleBlockType } from '@/interfaces/Article'

interface Props {
  initialUri?: string
}

const $props = defineProps<Props>()
const $emit = defineEmits<{
  saved: [uri: string]
  'dirty-change': [dirty: boolean]
}>()

const { t } = useI18n()
const toast = useToast()

const {
  article,
  uri,
  title,
  subtitle,
  published_at,
  blocks,
  loading,
  saving,
  error,
  conflict,
  dirty,
  open,
  save,
  reset
} = useArticleDraft()

const PREVIEW_DEBOUNCE_MS = 3000

const preview_uri = ref('')
const preview_title = ref('')
const preview_subtitle = ref('')
const preview_published_at = ref<string | null>(null)
const preview_blocks = ref<ArticleBlock[]>([])

let preview_timer: number | undefined

const error_message = computed(() => {
  if (error.value === 'save') return 'article.editor.saveFailed'
  if (error.value === 'create') return 'article.editor.conflict'

  return 'article.editor.loadFailed'
})

const publish_local = computed({
  get: () => (published_at.value ? dayjs(published_at.value).format('YYYY-MM-DDTHH:mm') : ''),
  set: (value: string) => {
    published_at.value = value ? dayjs(value).toISOString() : null
  }
})

const preview_article = computed<Article>(() => ({
  uri: preview_uri.value,
  author: article.value?.author || { email: '', display_name: '' },
  title: preview_title.value,
  subtitle: preview_subtitle.value,
  blocks: preview_blocks.value,
  summary: '',
  cover: '',
  created_at: article.value?.created_at || '',
  updated_at: article.value?.updated_at || '',
  published_at: preview_published_at.value
}))

function cloneBlocks(value: ArticleBlock[]): ArticleBlock[] {
  return JSON.parse(JSON.stringify(value || [])) as ArticleBlock[]
}

function refreshPreview(): void {
  window.clearTimeout(preview_timer)
  preview_timer = undefined

  preview_uri.value = article.value?.uri || ''
  preview_title.value = title.value
  preview_subtitle.value = subtitle.value
  preview_published_at.value = published_at.value
  preview_blocks.value = cloneBlocks(blocks.value)
}

function schedulePreview(): void {
  if (!article.value) return

  window.clearTimeout(preview_timer)
  preview_timer = window.setTimeout(refreshPreview, PREVIEW_DEBOUNCE_MS)
}

type DrawingEditor = {
  saveScene: () => Promise<string | null>
}

const drawing_refs = new Map<string, DrawingEditor>()

function setDrawingRef(id: string, instance: unknown): void {
  if (instance) {
    drawing_refs.set(id, instance as DrawingEditor)
  } else {
    drawing_refs.delete(id)
  }
}

const block_elements = new Map<string, HTMLElement>()

function setBlockElement(id: string, instance: unknown): void {
  if (instance) {
    const element = (instance as { $el?: unknown }).$el
    if (element instanceof HTMLElement) block_elements.set(id, element)
  } else {
    block_elements.delete(id)
  }
}

function addBlock(type: ArticleBlockType): void {
  blocks.value.push(createArticleBlock(type))
}

function removeBlock(index: number): void {
  const [removed] = blocks.value.splice(index, 1)
  if (removed?.type === 'excalidraw') drawing_refs.delete(removed.id)
}

function moveBlock(index: number, delta: number): void {
  const target = index + delta
  if (target < 0 || target >= blocks.value.length) return

  const [block] = blocks.value.splice(index, 1)
  if (!block) return

  blocks.value.splice(target, 0, block)

  nextTick(() => {
    block_elements.get(block.id)?.scrollIntoView({
      behavior: 'smooth',
      block: 'center'
    })
  })
}

async function saveArticle(): Promise<void> {
  if (!article.value || saving.value) return

  for (const block of blocks.value) {
    if (block.type !== 'excalidraw') continue

    const editor = drawing_refs.get(block.id)
    if (!editor) continue

    const scene_path = await editor.saveScene()
    if (!scene_path) {
      toast.error(t('article.editor.sceneFailed'))
      return
    }

    block.scene_path = scene_path
  }

  const saved = await save()

  if (saved) {
    $emit('saved', uri.value)
    toast.success(t('article.editor.saved'))
  } else if (conflict.value) {
    toast.error(t('article.editor.conflict'))
  } else {
    toast.error(t('article.editor.saveFailed'))
  }
}

function onKeyDown(event: KeyboardEvent): void {
  if (!(event.ctrlKey || event.metaKey) || event.key.toLowerCase() !== 's') return
  if (event.repeat) return

  event.preventDefault()
  void saveArticle()
}

function onBeforeUnload(event: BeforeUnloadEvent): void {
  if (!dirty.value) return

  event.preventDefault()
  event.returnValue = ''
}

onMounted(() => {
  window.addEventListener('keydown', onKeyDown)
  window.addEventListener('beforeunload', onBeforeUnload)
})

watch(
  article,
  (value) => {
    if (value) {
      refreshPreview()
      return
    }

    window.clearTimeout(preview_timer)
    preview_timer = undefined
    preview_uri.value = ''
    preview_title.value = ''
    preview_subtitle.value = ''
    preview_published_at.value = null
    preview_blocks.value = []
  },
  { immediate: true }
)

watch([title, subtitle, uri, published_at, blocks], schedulePreview, { deep: true })

watch(
  () => $props.initialUri,
  async (value) => {
    if (value) {
      await open(value)
    } else {
      reset()
    }
  },
  { immediate: true }
)

onUnmounted(() => {
  window.removeEventListener('keydown', onKeyDown)
  window.clearTimeout(preview_timer)
})
</script>
