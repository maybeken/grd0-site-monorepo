<template>
  <div class="flex h-full flex-col gap-5">
    <div class="flex flex-col gap-1">
      <p class="text-lg font-bold">{{ $t('article.editor.details') }}</p>
      <span class="text-sm text-secondary">{{ status_label }}</span>
    </div>

    <div class="flex flex-col gap-2">
      <div class="flex items-center justify-between gap-2">
        <label for="article-sidebar-uri">{{ $t('article.editor.uriLabel') }}</label>
        <button
          v-if="$props.article.published_at"
          class="shrink-0 rounded-full border border-dotted p-1 text-secondary"
          :title="copied ? $t('article.editor.copied') : $t('article.editor.copyLink')"
          :aria-label="copied ? $t('article.editor.copied') : $t('article.editor.copyLink')"
          @click="copyLink"
        >
          <Icon :icon="copied ? 'mynaui:check' : 'mynaui:copy'" height="1rem"></Icon>
        </button>
      </div>
      <div class="flex flex-wrap items-center gap-1 text-sm">
        <span class="text-secondary">/article/</span>
        <input
          id="article-sidebar-uri"
          class="min-w-0 flex-1 rounded-full border border-dotted bg-background px-3 py-1"
          type="text"
          :placeholder="$t('article.editor.uriPlaceholder')"
          v-model="uri"
        />
      </div>
      <p v-if="!valid" class="text-sm text-red-400">{{ $t('article.editor.invalidUri') }}</p>
    </div>

    <div class="flex flex-col gap-2">
      <label for="article-sidebar-publish">{{ $t('article.editor.publishDate') }}</label>
      <input
        id="article-sidebar-publish"
        class="rounded-xl border border-dotted bg-background px-2 py-1"
        type="datetime-local"
        v-model="publish_local"
      />
      <div class="flex flex-wrap gap-2">
        <button
          v-if="status !== 'published'"
          class="rounded-full border border-dotted px-3 py-1 text-sm disabled:opacity-50"
          :disabled="saving"
          @click="$emit('publish-now')"
        >
          {{ $t('article.editor.publishNow') }}
        </button>
        <button
          v-if="status !== 'draft'"
          class="rounded-full border border-dotted px-3 py-1 text-sm disabled:opacity-50"
          :disabled="saving"
          @click="$emit('unpublish')"
        >
          {{ $t('article.editor.unpublish') }}
        </button>
      </div>
    </div>

    <div class="flex flex-col gap-1 text-sm text-secondary">
      <span>{{ $t('article.editor.createdAt', { date: formatDate($props.article.created_at) }) }}</span>
      <span>{{ $t('article.editor.updatedAt', { date: formatDate($props.article.updated_at) }) }}</span>
    </div>

    <router-link
      :to="`/article/${$props.article.uri}`"
      target="_blank"
      class="rounded-full border border-dotted px-3 py-1 text-center text-sm"
    >
      {{ $t('article.editor.openPublic') }}
    </router-link>

    <button
      class="mt-auto rounded-full border border-red-500/40 px-3 py-1 text-sm text-red-400 disabled:opacity-50"
      :disabled="saving"
      @click="$emit('delete')"
    >
      {{ $t('common.delete') }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import dayjs from 'dayjs'
import LocalizedFormat from 'dayjs/plugin/localizedFormat'

import type { Article } from '@/interfaces/Article'

dayjs.extend(LocalizedFormat)

interface Props {
  article: Article
  saving: boolean
}

const $props = defineProps<Props>()

const $emit = defineEmits<{
  'publish-now': []
  unpublish: []
  delete: []
}>()

const uri = defineModel<string>('uri', { required: true })
const published_at = defineModel<string | null>('publishedAt', { required: true })

const { t } = useI18n()
const copied = ref(false)

let copy_timer: number | undefined

const valid = computed(() => /^[a-z0-9-]+$/.test(uri.value))

const publish_local = computed({
  get: () => (published_at.value ? dayjs(published_at.value).format('YYYY-MM-DDTHH:mm') : ''),
  set: (value: string) => {
    published_at.value = value ? dayjs(value).toISOString() : null
  }
})

const status = computed<'draft' | 'scheduled' | 'published'>(() => {
  if (!published_at.value) return 'draft'

  return dayjs(published_at.value).isAfter(dayjs()) ? 'scheduled' : 'published'
})

const status_keys = {
  draft: 'article.editor.statusDraft',
  scheduled: 'article.editor.statusScheduled',
  published: 'article.editor.statusPublished'
} as const

const status_label = computed(() => t(status_keys[status.value]))

function formatDate(value: string): string {
  return dayjs(value).format('LLL')
}

async function copyLink(): Promise<void> {
  try {
    await navigator.clipboard.writeText(
      `${window.location.origin}/article/${$props.article.uri}`
    )
    copied.value = true

    window.clearTimeout(copy_timer)
    copy_timer = window.setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch {
    copied.value = false
  }
}

onBeforeUnmount(() => {
  window.clearTimeout(copy_timer)
})
</script>
