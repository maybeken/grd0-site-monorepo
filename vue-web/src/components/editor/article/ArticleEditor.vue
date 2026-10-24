<template>
  <EditorPage class="lg:h-[calc(100vh-3rem)] lg:min-h-0">
    <EditorPageHeader
      :title="article?.title || $t('article.editor.title')"
      title-sr-only
      :back-to="back_to"
      :back-label="$t('editor.ui.backTo', { section: $t('editor.articles') })"
    >
      <template #actions>
        <span v-if="article" class="text-sm font-medium text-foreground/60">{{ status_label }}</span>
        <span v-if="article" class="text-sm font-medium text-foreground/60">
          {{ saving ? $t('common.saving') : dirty ? $t('article.editor.unsaved') : $t('common.saved') }}
        </span>
        <EditorButton variant="secondary" class="lg:hidden" @click="sidebar_open = !sidebar_open">
          {{ $t('article.editor.details') }}
        </EditorButton>
        <EditorButton
          variant="primary"
          :disabled="saving || loading || !dirty"
          :loading="saving"
          aria-keyshortcuts="Control+S Meta+S"
          @click="saveArticle()"
        >
          {{ $t('common.save') }}
        </EditorButton>
      </template>
    </EditorPageHeader>

    <EditorError v-if="conflict" :message="$t('article.editor.conflict')"></EditorError>
    <EditorError v-else-if="error" :message="$t(error_message)"></EditorError>

    <div v-if="loading" class="grid grid-cols-1 gap-4 lg:grid-cols-3">
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
      <div class="flex flex-col gap-3 rounded-xl border border-dotted p-4">
        <Skeleton h="md" w="1/2" :loading="true"></Skeleton>
        <Skeleton h="lg" w="full" :loading="true"></Skeleton>
      </div>
    </div>

    <div
      v-else-if="article"
      class="flex flex-col gap-4 lg:min-h-0 lg:flex-1 lg:flex-row lg:items-stretch"
    >
      <div class="flex min-w-0 flex-1 flex-col gap-4 lg:min-h-0">
        <div class="flex flex-col gap-3">
          <input
            class="w-full border-b border-dotted bg-transparent px-1 py-2 text-3xl outline-none focus:outline focus:outline-2 focus:outline-offset-2 focus:outline-foreground/60"
            type="text"
            :placeholder="$t('article.editor.titlePlaceholder')"
            v-model="title"
          />
          <input
            class="w-full border-b border-dotted bg-transparent px-1 py-2 text-2xl outline-none focus:outline focus:outline-2 focus:outline-offset-2 focus:outline-foreground/60"
            type="text"
            :placeholder="$t('article.editor.subtitlePlaceholder')"
            v-model="subtitle"
          />
        </div>

        <div class="flex gap-2 lg:hidden">
          <button
            class="rounded-full border border-dotted px-3 py-1"
            :class="view === 'edit' ? 'bg-shade' : ''"
            @click="view = 'edit'"
          >
            {{ $t('article.editor.editTab') }}
          </button>
          <button
            class="rounded-full border border-dotted px-3 py-1"
            :class="view === 'preview' ? 'bg-shade' : ''"
            @click="view = 'preview'"
          >
            {{ $t('article.editor.previewTab') }}
          </button>
        </div>

        <div class="flex flex-col gap-4 lg:min-h-0 lg:flex-1 lg:grid lg:grid-cols-2">
          <section
            ref="editor_pane"
            class="min-h-0 flex-col gap-6 overflow-y-auto p-4 max-lg:h-[calc(100vh-19rem)] lg:flex"
            :class="view === 'edit' ? 'flex' : 'hidden'"
          >
            <div class="flex flex-col gap-6">
              <div
                v-for="(block, index) in blocks"
                :key="block.id"
                :ref="(element) => setBlockElement(block.id, element)"
                :data-block-id="block.id"
              >
                <BlockItem
                  :type="block.type"
                  :index="index"
                  :count="blocks.length"
                  :active="active_block_id === block.id"
                  @move-up="moveBlock(index, -1)"
                  @move-down="moveBlock(index, 1)"
                  @duplicate="duplicateBlock(index)"
                  @remove="removeBlock(index)"
                >
                  <MarkdownBlockEditor
                    v-if="block.type === 'markdown'"
                    :block="block"
                  ></MarkdownBlockEditor>
                  <ImageBlockEditor
                    v-else-if="block.type === 'image'"
                    :block="block"
                  ></ImageBlockEditor>
                  <MusicBlockEditor
                    v-else-if="block.type === 'music'"
                    :block="block"
                  ></MusicBlockEditor>
                  <ExcalidrawBlockEditor
                    v-else-if="block.type === 'excalidraw'"
                    :ref="(element) => setDrawingRef(block.id, element)"
                    :block="block"
                    :uri="article.uri"
                    @request-save="saveArticle"
                  ></ExcalidrawBlockEditor>
                </BlockItem>
              </div>

              <div class="flex flex-col gap-2">
                <p class="text-lg font-bold">{{ $t('article.editor.addBlock') }}</p>
                <BlockPalette @add="addBlock"></BlockPalette>
              </div>
            </div>
          </section>

          <section
            ref="preview_pane"
            class="min-h-0 flex-col gap-4 overflow-y-auto p-4 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden max-lg:h-[calc(100vh-19rem)] lg:flex lg:border-l lg:border-dotted"
            :class="view === 'preview' ? 'flex' : 'hidden'"
          >
            <div class="flex flex-col gap-4">
              <BlockRenderer
                :blocks="preview_blocks"
                :uri="preview_uri"
                :version="preview_version"
                :draft="true"
                :eager="true"
              ></BlockRenderer>
            </div>
          </section>
        </div>
      </div>

      <aside
        class="min-h-0 flex-col gap-4 overflow-y-auto rounded-xl border border-dotted p-4 max-lg:max-h-[70vh] lg:flex lg:w-72 lg:shrink-0"
        :class="sidebar_open ? 'flex' : 'hidden'"
      >
        <ArticleSidebar
          v-model:uri="uri"
          v-model:published-at="published_at"
          :article="article"
          :saving="saving"
          @publish-now="publishNow"
          @unpublish="unpublish"
          @delete="purgeArticle"
        ></ArticleSidebar>
      </aside>
    </div>

    <div v-else class="flex flex-col items-start gap-3">
      <EditorError :message="$t('article.editor.loadFailed')">
        <template #actions>
          <EditorButton variant="secondary" :to="back_to">
            {{ $t('editor.ui.backTo', { section: $t('editor.articles') }) }}
          </EditorButton>
        </template>
      </EditorError>
    </div>
  </EditorPage>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import dayjs from 'dayjs'

import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorError from '@/components/editor/ui/EditorError.vue'
import EditorPage from '@/components/editor/ui/EditorPage.vue'
import EditorPageHeader from '@/components/editor/ui/EditorPageHeader.vue'
import { articleSceneId, createArticleBlock } from '@/helpers/article'
import { useSaveShortcut } from '@/composables/useSaveShortcut'
import { useToast } from '@/composables/useToast'
import { useUnsavedGuard } from '@/composables/useUnsavedGuard'
import {
  getArticleDraftSceneRaw,
  purgeArticleRaw,
  putArticleSceneRaw
} from '@/services/article'
import ArticleSidebar from './ArticleSidebar.vue'
import BlockItem from './BlockItem.vue'
import BlockPalette from './BlockPalette.vue'
import MarkdownBlockEditor from './blocks/MarkdownBlockEditor.vue'
import ImageBlockEditor from './blocks/ImageBlockEditor.vue'
import MusicBlockEditor from './blocks/MusicBlockEditor.vue'
import ExcalidrawBlockEditor from './blocks/ExcalidrawBlockEditor.vue'
import { useArticleDraft } from './useArticleDraft'
import { useScrollSync } from './useScrollSync'

import type { ArticleBlock, ArticleBlockType } from '@/interfaces/Article'

interface Props {
  uri: string
}

const $props = defineProps<Props>()
const $emit = defineEmits<{
  saved: [uri: string]
}>()

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
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

useUnsavedGuard(dirty)
useSaveShortcut(saveArticle, { enabled: () => dirty.value && !saving.value && !loading.value })

const PREVIEW_DEBOUNCE_MS = 750

const editor_pane = ref<HTMLElement | null>(null)
const preview_pane = ref<HTMLElement | null>(null)
const is_wide = ref(false)
const view = ref<'edit' | 'preview'>('edit')
const sidebar_open = ref(false)

const scroll_sync = useScrollSync(editor_pane, preview_pane, is_wide)
const { active_block_id } = scroll_sync

const preview_uri = ref('')
const preview_version = ref('')
const preview_blocks = ref<ArticleBlock[]>([])

let preview_timer: number | undefined
let media_query: MediaQueryList | null = null

const error_message = computed(() => {
  if (error.value === 'save') return 'article.editor.saveFailed'
  if (error.value === 'create') return 'article.editor.conflict'

  return 'article.editor.loadFailed'
})

const status_label = computed(() => {
  if (!published_at.value) return t('article.editor.statusDraft')

  return dayjs(published_at.value).isAfter(dayjs())
    ? t('article.editor.statusScheduled')
    : t('article.editor.statusPublished')
})

const back_to = computed(() => {
  const value = Number(route.query.page)
  if (!Number.isInteger(value) || value <= 0) return '/editor/articles'

  return { path: '/editor/articles', query: { page: String(value) } }
})

function cloneBlocks(value: ArticleBlock[]): ArticleBlock[] {
  return JSON.parse(JSON.stringify(value || [])) as ArticleBlock[]
}

function refreshPreview(): void {
  window.clearTimeout(preview_timer)
  preview_timer = undefined

  preview_uri.value = article.value?.uri || ''
  preview_version.value = article.value?.updated_at || ''
  preview_blocks.value = cloneBlocks(blocks.value)

  nextTick(() => scroll_sync.realign())
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
  if (instance instanceof HTMLElement) {
    block_elements.set(id, instance)
  } else {
    block_elements.delete(id)
  }
}

function addBlock(type: ArticleBlockType): void {
  blocks.value.push(createArticleBlock(type))
  nextTick(() => scroll_sync.realign())
}

function removeBlock(index: number): void {
  const [removed] = blocks.value.splice(index, 1)
  if (removed?.type === 'excalidraw') drawing_refs.delete(removed.id)

  nextTick(() => scroll_sync.realign())
}

async function duplicateBlock(index: number): Promise<void> {
  const block = blocks.value[index]
  if (!block || !article.value) return

  const id = crypto.randomUUID()
  let copy: ArticleBlock = { ...block, id }

  if (block.type === 'excalidraw' && block.scene_path) {
    const scene_id = articleSceneId(block.scene_path)

    if (scene_id) {
      try {
        const scene = await getArticleDraftSceneRaw(article.value.uri, scene_id).send()
        const response = await putArticleSceneRaw(
          article.value.uri,
          crypto.randomUUID(),
          scene as object
        ).send()

        if (!response?.scene_path) throw new Error('Missing scene path')

        copy = { ...block, id, scene_path: response.scene_path }
      } catch {
        toast.error(t('article.editor.sceneFailed'))
        return
      }
    }
  }

  const source_index = blocks.value.findIndex((item) => item.id === block.id)
  if (source_index === -1) return

  blocks.value.splice(source_index + 1, 0, copy)

  nextTick(() => {
    block_elements.get(copy.id)?.scrollIntoView({
      behavior: 'smooth',
      block: 'center'
    })
    scroll_sync.realign()
  })
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
    scroll_sync.realign()
  })
}

async function saveArticle(success_message = 'article.editor.saved'): Promise<void> {
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
    toast.success(t(success_message))
    await nextTick()
    scroll_sync.realign()
  } else if (conflict.value) {
    toast.error(t('article.editor.conflict'))
  } else {
    toast.error(t('article.editor.saveFailed'))
  }
}

async function publishNow(): Promise<void> {
  if (!article.value || saving.value) return

  published_at.value = dayjs().toISOString()
  await saveArticle('article.editor.published')
}

async function unpublish(): Promise<void> {
  if (!article.value || saving.value) return

  published_at.value = null
  await saveArticle('article.editor.unpublishedToast')
}

async function purgeArticle(): Promise<void> {
  if (!article.value || saving.value) return
  if (!window.confirm(t('article.editor.confirmDelete'))) return

  try {
    await purgeArticleRaw(article.value.uri).send()
    toast.success(t('article.editor.deleted'))
    reset()
    await router.push(back_to.value)
  } catch {
    toast.error(t('article.editor.deleteFailed'))
  }
}

function updateWidth(): void {
  is_wide.value = media_query?.matches ?? false
}

onMounted(() => {
  media_query = window.matchMedia('(min-width: 1024px)')
  updateWidth()
  media_query.addEventListener('change', updateWidth)
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
    preview_version.value = ''
    preview_blocks.value = []
    scroll_sync.reset()
  },
  { immediate: true }
)

watch(blocks, schedulePreview, { deep: true })

watch(
  view,
  (value) => {
    if (value !== 'preview') return

    nextTick(() => {
      if (is_wide.value) {
        scroll_sync.realign()
      } else {
        scroll_sync.scrollPreviewToActive()
      }
    })
  }
)

watch(
  () => $props.uri,
  async (value) => {
    if (!value) {
      reset()
      return
    }

    if (article.value?.uri === value) return

    await open(value)
  },
  { immediate: true }
)

onUnmounted(() => {
  media_query?.removeEventListener('change', updateWidth)
  window.clearTimeout(preview_timer)
})
</script>
