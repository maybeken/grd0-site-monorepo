<template>
  <div class="flex h-full flex-col gap-5">
    <div class="flex flex-col gap-1">
      <h2 class="text-xl font-bold">{{ $t('article.editor.details') }}</h2>
      <span class="text-sm font-medium text-foreground/60">{{ status_label }}</span>
    </div>

    <div class="flex flex-col gap-2">
      <div class="flex items-center justify-between gap-2">
        <label for="article-sidebar-uri" class="text-lg font-medium">
          {{ $t('article.editor.uriLabel') }}
        </label>
        <EditorButton
          v-if="$props.article.published_at"
          variant="icon"
          :label="copied ? $t('article.editor.copied') : $t('article.editor.copyLink')"
          @click="copyLink"
        >
          <Icon :icon="copied ? 'mynaui:check' : 'mynaui:copy'" height="1rem"></Icon>
        </EditorButton>
      </div>
      <p class="text-sm font-medium text-foreground/60">
        {{ $t('article.editor.uriPrefix') }}
      </p>
      <EditorInput
        id="article-sidebar-uri"
        v-model="uri"
        :placeholder="$t('article.editor.uriPlaceholder')"
      ></EditorInput>
      <p v-if="!valid" class="text-sm font-medium text-red-400">
        {{ $t('article.editor.invalidUri') }}
      </p>
    </div>

    <div class="flex flex-col gap-2">
      <label for="article-sidebar-publish" class="text-lg font-medium">
        {{ $t('article.editor.publishDate') }}
      </label>
      <EditorInput
        id="article-sidebar-publish"
        v-model="publish_local"
        type="datetime-local"
      ></EditorInput>
      <div class="flex flex-wrap gap-2">
        <EditorButton
          v-if="status !== 'published'"
          variant="secondary"
          size="sm"
          :disabled="saving"
          @click="$emit('publish-now')"
        >
          {{ $t('article.editor.publishNow') }}
        </EditorButton>
        <EditorButton
          v-if="status !== 'draft'"
          variant="secondary"
          size="sm"
          :disabled="saving"
          @click="$emit('unpublish')"
        >
          {{ $t('common.unpublish') }}
        </EditorButton>
      </div>
    </div>

    <div class="flex flex-col gap-1 text-sm font-medium text-foreground/60">
      <span>{{ $t('article.editor.createdAt', { date: formatDate($props.article.created_at) }) }}</span>
      <span>{{ $t('article.editor.updatedAt', { date: formatDate($props.article.updated_at) }) }}</span>
    </div>

    <EditorButton variant="secondary" :to="`/article/${$props.article.uri}`" target="_blank">
      {{ $t('article.editor.openPublic') }}
    </EditorButton>

    <EditorButton
      variant="destructive"
      class="mt-auto"
      :disabled="saving"
      @click="$emit('delete')"
    >
      {{ $t('common.delete') }}
    </EditorButton>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import dayjs from 'dayjs'
import LocalizedFormat from 'dayjs/plugin/localizedFormat'

import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorInput from '@/components/editor/ui/EditorInput.vue'
import Icon from '@/components/shared/Icon.vue'

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
      `${window.location.origin}${t('article.editor.uriPrefix')}${$props.article.uri}`
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
