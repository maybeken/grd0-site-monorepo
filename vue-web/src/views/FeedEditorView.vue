<template>
  <div class="flex flex-col gap-6 px-4 md:px-16">
    <div class="flex justify-between items-center">
      <p class="text-2xl font-bold">{{ $t('feed.editor.title') }}<CursorBlink /></p>
      <div class="flex gap-2">
        <button
          class="px-4 py-2 border border-dotted rounded-full hover:bg-accent/20 disabled:brightness-50"
          :disabled="saving"
          @click="saveDraft"
        >
          {{ $t('feed.editor.saveDraft') }}
        </button>
        <button
          class="px-4 py-2 border border-dotted rounded-full bg-accent/20 hover:bg-accent/30 disabled:brightness-50"
          :disabled="saving"
          @click="publish"
        >
          {{ $t('feed.editor.publish') }}
        </button>
      </div>
    </div>

    <div v-if="message" class="px-4 py-2 rounded-xl" :class="messageType === 'error' ? 'bg-red-900/30 text-red-300' : 'bg-emerald-900/30 text-emerald-300'">
      {{ message }}
    </div>

    <input
      v-model="title"
      class="text-3xl w-full rounded-full px-4 py-2 bg-shade"
      type="text"
      :placeholder="$t('feed.editor.titlePlaceholder')"
    />
    <input
      v-model="subtitle"
      class="text-xl w-full rounded-full px-4 py-2 bg-shade"
      type="text"
      :placeholder="$t('feed.editor.subtitlePlaceholder')"
    />

    <div class="flex gap-2">
      <button
        class="px-4 py-2 border border-dotted rounded-full hover:bg-accent/20"
        @click="showAddMenu = !showAddMenu"
      >
        {{ $t('feed.editor.addBlock') }}
      </button>
      <div v-if="showAddMenu" class="relative">
        <div class="absolute z-10 flex flex-col bg-shade rounded-xl border border-accent shadow-md">
          <button class="px-4 py-2 hover:bg-accent/20 text-left" @click="addBlogBlock">
            {{ $t('feed.editor.addBlog') }}
          </button>
          <button class="px-4 py-2 hover:bg-accent/20 text-left" @click="addMusicBlock">
            {{ $t('feed.editor.addMusic') }}
          </button>
          <button class="px-4 py-2 hover:bg-accent/20 text-left" @click="addReferenceBlock">
            {{ $t('feed.editor.addReference') }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="validationErrors.length" class="flex flex-col gap-1 px-4 py-2 rounded-xl bg-red-900/30 text-red-300 text-sm">
      <p v-for="err in validationErrors" :key="err">{{ err }}</p>
    </div>

    <div v-for="(block, idx) in blocks" :key="block.id" class="flex flex-col gap-2 bg-shade rounded-2xl p-4" :class="{ 'ring-1 ring-accent': selectedBlockIdx === idx }">
      <div class="flex justify-between items-center cursor-pointer" @click="selectedBlockIdx = selectedBlockIdx === idx ? -1 : idx">
        <span class="text-sm text-secondary font-bold">
          #{{ idx + 1 }} - {{ block.type === 'blog' ? (block.mode === 'reference' ? $t('feed.editor.refBlog') : $t('feed.editor.newBlog')) : (block.mode === 'reference' ? $t('feed.editor.refMusic') : $t('feed.editor.newMusic')) }}
        </span>
        <div class="flex gap-1">
          <button
            class="p-1 border border-dotted rounded-full hover:bg-accent/20 disabled:brightness-50"
            :disabled="idx === 0"
            @click="moveBlock(idx, -1)"
            :title="$t('feed.editor.moveUp')"
          >
            <Icon icon="mynaui:chevron-up" height="1rem"></Icon>
          </button>
          <button
            class="p-1 border border-dotted rounded-full hover:bg-accent/20 disabled:brightness-50"
            :disabled="idx === blocks.length - 1"
            @click="moveBlock(idx, 1)"
            :title="$t('feed.editor.moveDown')"
          >
            <Icon icon="mynaui:chevron-down" height="1rem"></Icon>
          </button>
          <button
            class="p-1 border border-dotted rounded-full hover:bg-red-500/20 text-red-400"
            @click="removeBlock(idx)"
            :title="$t('common.delete')"
          >
            <Icon icon="mynaui:trash" height="1rem"></Icon>
          </button>
        </div>
      </div>

      <template v-if="block.type === 'blog' && block.mode === 'create'">
        <input
          v-model="block.createData.title"
          class="text-xl w-full rounded-full px-4 py-1 bg-background"
          type="text"
          :placeholder="$t('feed.editor.blogTitlePlaceholder')"
        />
        <div class="flex gap-4">
          <div class="w-1/2">
            <textarea
              v-model="block.createData.content"
              class="w-full h-64 rounded-xl px-4 py-2 bg-background font-mono text-sm"
              :placeholder="$t('feed.editor.markdownPlaceholder')"
            ></textarea>
          </div>
          <div class="w-1/2 max-h-64 overflow-y-auto rounded-xl px-4 py-2 bg-background">
            <MarkdownDisplay :md="block.createData.content ?? ''" :lazy_loading="false"></MarkdownDisplay>
          </div>
        </div>
      </template>

      <template v-else-if="block.type === 'music' && block.mode === 'create'">
        <div class="flex flex-col gap-2">
          <input
            v-model="block.createData.v"
            class="w-full rounded-full px-4 py-1 bg-background"
            type="text"
            :placeholder="$t('feed.editor.youtubeIdPlaceholder')"
          />
          <input
            v-model="block.createData.title"
            class="w-full rounded-full px-4 py-1 bg-background"
            type="text"
            :placeholder="$t('feed.editor.musicTitlePlaceholder')"
          />
          <input
            v-model="block.createData.artist"
            class="w-full rounded-full px-4 py-1 bg-background"
            type="text"
            :placeholder="$t('feed.editor.artistPlaceholder')"
          />
        </div>
      </template>

      <template v-else-if="block.mode === 'reference'">
        <div class="flex flex-col gap-2">
          <div class="flex gap-2">
            <button
              class="px-3 py-1 rounded-full text-sm border border-dotted"
              :class="block.refType === 'blog' ? 'bg-accent/30' : ''"
              @click="block.refType = 'blog'; block.refId = ''"
            >
              {{ $t('feed.editor.searchBlog') }}
            </button>
            <button
              class="px-3 py-1 rounded-full text-sm border border-dotted"
              :class="block.refType === 'music' ? 'bg-accent/30' : ''"
              @click="block.refType = 'music'; block.refId = ''"
            >
              {{ $t('feed.editor.searchMusic') }}
            </button>
          </div>
          <input
            v-model="block.searchQuery"
            class="w-full rounded-full px-4 py-1 bg-background"
            type="text"
            :placeholder="$t('common.search')"
            @input="searchReferences(block)"
          />
          <div v-if="block.searchResults?.length" class="flex flex-col gap-1 max-h-48 overflow-y-auto">
            <button
              v-for="result in block.searchResults"
              :key="result.id"
              class="px-3 py-1 text-left rounded-lg hover:bg-accent/20 text-sm"
              @click="selectReference(block, result)"
            >
              {{ result.label }}
            </button>
          </div>
          <p v-if="block.refId" class="text-sm text-emerald-400">
            {{ $t('feed.editor.selected') }}: {{ block.refLabel }}
          </p>
        </div>
      </template>
    </div>

    <div
      v-if="selectedBlockIdx >= 0 && blocks[selectedBlockIdx]?.type === 'blog' && blocks[selectedBlockIdx]?.mode === 'reference' && blocks[selectedBlockIdx]?.refId"
      class="bg-shade rounded-2xl p-4"
    >
      <p class="text-sm text-secondary font-bold mb-2">{{ $t('feed.editor.editingBlog') }}</p>
      <SideBySideEditor :key="editorKey" :uri="blocks[selectedBlockIdx]!.refId!" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import dayjs from 'dayjs'

import { getFeedBySlugRaw, upsertFeedRaw } from '@/services/feed'
import { listBlogPostRaw } from '@/services/blogPost'
import { dataInstance } from '@/services/api'
import SideBySideEditor from '@/components/editor/SideBySideEditor.vue'
import type { FeedEntry, FeedBlockInput, EditorBlock, FeedBlockCreateData } from '@/interfaces/Feed'
import type { BlogPost } from '@/interfaces/Blog'
import type { Music } from '@/interfaces/Music'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const existingSlug = typeof route.params.slug === 'string' ? route.params.slug : route.params.slug?.[0]

const title = ref('')
const subtitle = ref('')
const blocks = reactive<EditorBlock[]>([])
const showAddMenu = ref(false)
const saving = ref(false)
const message = ref('')
const messageType = ref<'success' | 'error'>('success')
const validationErrors = ref<string[]>([])
const selectedBlockIdx = ref(-1)

const editorKey = computed(() => {
  if (selectedBlockIdx.value < 0) return 'none'
  const block = blocks[selectedBlockIdx.value]
  return block ? `${block.id}-${block.refId || 'new'}` : 'none'
})

let blockIdCounter = 0
function nextBlockId() {
  return 'block-' + (++blockIdCounter)
}

function addBlogBlock() {
  blocks.push({
    id: nextBlockId(),
    type: 'blog',
    mode: 'create',
    order: blocks.length,
    createData: { title: '', content: '' }
  })
  showAddMenu.value = false
}

function addMusicBlock() {
  blocks.push({
    id: nextBlockId(),
    type: 'music',
    mode: 'create',
    order: blocks.length,
    createData: { v: '', title: '', artist: '' }
  })
  showAddMenu.value = false
}

function addReferenceBlock() {
  blocks.push({
    id: nextBlockId(),
    type: 'blog',
    mode: 'reference',
    order: blocks.length,
    createData: {},
    refType: 'blog',
    searchQuery: '',
    searchResults: [],
    refId: '',
    refLabel: ''
  })
  showAddMenu.value = false
}

function removeBlock(idx: number) {
  blocks.splice(idx, 1)
  reindexBlocks()
  selectedBlockIdx.value = -1
}

function moveBlock(idx: number, direction: number) {
  const newIdx = idx + direction
  if (newIdx < 0 || newIdx >= blocks.length) return
  const current = blocks[idx]!
  const target = blocks[newIdx]!
  blocks[idx] = target
  blocks[newIdx] = current
  reindexBlocks()
  if (selectedBlockIdx.value === idx) selectedBlockIdx.value = newIdx
  else if (selectedBlockIdx.value === newIdx) selectedBlockIdx.value = idx
  else selectedBlockIdx.value = -1
}

function reindexBlocks() {
  blocks.forEach((b, i) => (b.order = i))
}

async function searchReferences(block: any) {
  const query = block.searchQuery
  if (!query || query.length < 2) {
    block.searchResults = []
    return
  }

  try {
    if (block.refType === 'blog') {
      const posts = (await listBlogPostRaw().send()) as unknown as BlogPost[]
      block.searchResults = posts
        .filter((p) => p.title.toLowerCase().includes(query.toLowerCase()) || p.uri.toLowerCase().includes(query.toLowerCase()))
        .map((p) => ({ id: p.uri, label: p.title, ref: p.uri }))
    } else {
      const music = (await dataInstance.Get<Music[]>('/music').send()) as unknown as Music[]
      block.searchResults = music
        .filter((m) => m.title.toLowerCase().includes(query.toLowerCase()) || (m.artist ?? '').toLowerCase().includes(query.toLowerCase()))
        .map((m) => ({ id: m.v, label: `${m.title} - ${m.artist ?? ''}`, ref: m.v }))
    }
  } catch {
    block.searchResults = []
  }
}

function selectReference(block: any, result: any) {
  block.refId = result.ref
  block.refLabel = result.label
  block.type = block.refType
  block.searchResults = []
  block.searchQuery = ''
}

function validate(): boolean {
  validationErrors.value = []

  if (!title.value.trim()) {
    validationErrors.value.push(t('feed.editor.validationTitleRequired'))
  }
  if (blocks.length === 0) {
    validationErrors.value.push(t('feed.editor.validationBlocksRequired'))
  }

  blocks.forEach((block, idx) => {
    if (block.mode === 'create' && block.type === 'blog' && !block.createData?.content?.trim()) {
      validationErrors.value.push(t('feed.editor.validationBlogContentRequired', { index: idx + 1 }))
    }
    if (block.mode === 'create' && block.type === 'music' && !block.createData?.v?.trim()) {
      validationErrors.value.push(t('feed.editor.validationMusicVidRequired', { index: idx + 1 }))
    }
  })

  return validationErrors.value.length === 0
}

function buildPayload(published: boolean): { slug: string; payload: any } {
  const slug = existingSlug || title.value.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '')

  const blockInputs: FeedBlockInput[] = blocks.map((block, idx) => {
    if (block.mode === 'reference') {
      return {
        type: block.type as 'blog' | 'music',
        order: idx,
        ref: block.refId
      }
    }

    return {
      type: block.type as 'blog' | 'music',
      order: idx,
      create_data: block.createData
    }
  })

  return {
    slug,
    payload: {
      title: title.value,
      subtitle: subtitle.value,
      published_at: published ? dayjs().toISOString() : null,
      blocks: blockInputs
    }
  }
}

async function saveDraft() {
  if (!validate()) return
  saving.value = true
  message.value = ''

  const { slug, payload } = buildPayload(false)

  try {
    await upsertFeedRaw(slug, payload).send()
    message.value = t('feed.editor.draftSaved')
    messageType.value = 'success'

    if (!existingSlug) {
      router.replace(`/feed/editor/${slug}`)
    }
  } catch {
    message.value = t('feed.editor.saveFailed')
    messageType.value = 'error'
  } finally {
    saving.value = false
  }
}

async function publish() {
  if (!validate()) return
  saving.value = true
  message.value = ''

  const { slug, payload } = buildPayload(true)

  try {
    await upsertFeedRaw(slug, payload).send()
    message.value = t('feed.editor.published')
    messageType.value = 'success'

    if (!existingSlug) {
      router.replace(`/feed/editor/${slug}`)
    }
  } catch {
    message.value = t('feed.editor.saveFailed')
    messageType.value = 'error'
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  if (!existingSlug) return

  try {
    const entry = (await getFeedBySlugRaw(existingSlug).send()) as unknown as FeedEntry
    title.value = entry.title
    subtitle.value = entry.subtitle ?? ''

    entry.blocks.forEach((block) => {
      if (block.blog) {
        blocks.push({
          id: nextBlockId(),
          type: 'blog',
          mode: 'reference',
          order: block.order,
          createData: {},
          refType: 'blog',
          refId: block.blog.uri,
          refLabel: block.blog.title,
          searchQuery: '',
          searchResults: []
        })
      } else if (block.music) {
        blocks.push({
          id: nextBlockId(),
          type: 'music',
          mode: 'reference',
          order: block.order,
          createData: {},
          refType: 'music',
          refId: block.music.v,
          refLabel: `${block.music.title} - ${block.music.artist ?? ''}`,
          searchQuery: '',
          searchResults: []
        })
      }
    })
  } catch {
    message.value = t('feed.editor.loadFailed')
    messageType.value = 'error'
  }
})
</script>
