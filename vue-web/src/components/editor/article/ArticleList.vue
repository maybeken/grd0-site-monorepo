<template>
  <div class="flex flex-col gap-3">
    <div class="flex flex-wrap items-center gap-2">
      <p class="text-lg font-bold">{{ $t('article.editor.listTitle') }}</p>
      <button
        class="ml-auto rounded-full border border-dotted px-3 py-1"
        @click="creating = !creating"
      >
        {{ $t('article.editor.newArticle') }}
      </button>
    </div>

    <div v-if="creating" class="flex flex-wrap items-center gap-2">
      <input
        class="rounded-full px-3 py-1 bg-background border border-dotted"
        type="text"
        :placeholder="$t('article.editor.uriPlaceholder')"
        v-model="new_uri"
      />
      <button
        class="rounded-full border border-dotted px-3 py-1 disabled:opacity-50"
        :disabled="busy"
        @click="create"
      >
        {{ $t('article.editor.create') }}
      </button>
      <button class="rounded-full border border-dotted px-3 py-1" @click="creating = false">
        {{ $t('common.cancel') }}
      </button>
    </div>

    <p v-if="conflict" class="text-red-400">{{ $t('article.editor.conflict') }}</p>
    <p v-else-if="invalid" class="text-red-400">{{ $t('article.editor.invalidUri') }}</p>
    <p v-else-if="failed" class="text-red-400">{{ $t('article.editor.loadFailed') }}</p>

    <div v-if="loading" class="flex flex-col gap-2">
      <Skeleton h="sm" w="full" :loading="true"></Skeleton>
      <Skeleton h="sm" w="full" :loading="true"></Skeleton>
      <Skeleton h="sm" w="full" :loading="true"></Skeleton>
    </div>

    <p v-else-if="entries.length === 0" class="text-secondary">
      {{ $t('article.editor.selectArticle') }}
    </p>

    <table v-else class="w-full table-auto">
      <thead class="border-b">
        <tr>
          <th class="py-2 text-left">{{ $t('article.editor.columnTitle') }}</th>
          <th class="py-2 text-left">{{ $t('article.editor.columnUri') }}</th>
          <th class="py-2 text-left">{{ $t('article.editor.columnPublished') }}</th>
          <th class="py-2">{{ $t('common.actions') }}</th>
        </tr>
      </thead>
      <tbody class="border-b">
        <tr
          v-for="entry in entries"
          :key="entry.uri"
          :class="entry.uri === $props.selected ? 'bg-shade' : ''"
        >
          <td class="py-2">{{ entry.title || entry.uri }}</td>
          <td class="py-2 text-secondary">/article/{{ entry.uri }}</td>
          <td class="py-2 text-secondary">{{ formatPublished(entry.published_at) }}</td>
          <td class="flex justify-center gap-2 py-2">
            <button
              class="rounded-xl border border-dotted px-4 py-px"
              @click="$emit('open', entry.uri)"
            >
              {{ $t('article.editor.open') }}
            </button>
            <button
              class="rounded-xl border border-dotted px-4 py-px"
              @click="unpublish(entry.uri)"
            >
              {{ $t('common.delete') }}
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { invalidateCache } from 'alova'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'

import { isArticle } from '@/helpers/article'
import {
  createArticleRaw,
  listArticlesAdminRaw,
  unpublishArticleRaw
} from '@/services/article'

import type { ArticleListEntry } from '@/interfaces/Article'

dayjs.extend(relativeTime)

interface Props {
  selected?: string
}

const $props = defineProps<Props>()
const $emit = defineEmits<{
  open: [uri: string]
  unpublished: [uri: string]
}>()

const { t } = useI18n()

const entries = ref<ArticleListEntry[]>([])
const loading = ref(true)
const busy = ref(false)
const creating = ref(false)
const new_uri = ref('')
const conflict = ref(false)
const invalid = ref(false)
const failed = ref(false)

async function refresh(): Promise<void> {
  loading.value = true

  try {
    await invalidateCache()
    const response = await listArticlesAdminRaw().send()
    entries.value = Array.isArray(response) ? response : []
  } catch {
    entries.value = []
  } finally {
    loading.value = false
  }
}

function formatPublished(value: string | null): string {
  if (!value) return t('article.editor.unpublished')

  const date = dayjs(value)
  if (date.isBefore(dayjs())) return t('article.editor.publishedAgo', { time: date.fromNow() })

  return t('article.editor.scheduledAgo', { time: date.fromNow() })
}

async function create(): Promise<void> {
  const target = new_uri.value.trim()

  conflict.value = false
  invalid.value = false
  failed.value = false

  if (!/^[a-z0-9-]+$/.test(target)) {
    invalid.value = true
    return
  }

  if (!confirm(t('article.editor.confirmCreate', { uri: target }))) return

  busy.value = true

  try {
    const response = await createArticleRaw({ uri: target }).send()

    if (isArticle(response)) {
      new_uri.value = ''
      creating.value = false
      await refresh()
      $emit('open', target)
      return
    }

    await refresh()
    conflict.value = entries.value.some((entry) => entry.uri === target)
    if (!conflict.value) failed.value = true
  } catch {
    await refresh()
    conflict.value = entries.value.some((entry) => entry.uri === target)
    if (!conflict.value) failed.value = true
  } finally {
    busy.value = false
  }
}

async function unpublish(uri: string): Promise<void> {
  if (!confirm(t('article.editor.confirmUnpublish'))) return

  busy.value = true
  conflict.value = false
  invalid.value = false
  failed.value = false

  try {
    await unpublishArticleRaw(uri).send()
    await refresh()
    $emit('unpublished', uri)
  } catch {
    failed.value = true
  } finally {
    busy.value = false
  }
}

onMounted(refresh)

defineExpose({ refresh })
</script>
