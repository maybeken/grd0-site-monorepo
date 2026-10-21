<template>
  <div class="flex flex-col gap-3">
    <div class="flex flex-wrap items-center gap-2">
      <p class="text-lg font-bold">{{ $t('article.editor.listTitle') }}</p>
      <button
        class="ml-auto rounded-full border border-dotted px-3 py-1 disabled:opacity-50"
        :disabled="busy"
        @click="create"
      >
        {{ $t('article.editor.newArticle') }}
      </button>
    </div>

    <p v-if="conflict" class="text-red-400">{{ $t('article.editor.conflict') }}</p>
    <p v-else-if="create_failed" class="text-red-400">
      {{ $t('article.editor.createFailed') }}
    </p>

    <div v-if="loading" class="flex flex-col gap-2">
      <Skeleton h="sm" w="full" :loading="true"></Skeleton>
      <Skeleton h="sm" w="full" :loading="true"></Skeleton>
      <Skeleton h="sm" w="full" :loading="true"></Skeleton>
    </div>

    <p v-else-if="list_failed" class="text-red-400">{{ $t('article.editor.listFailed') }}</p>

    <p v-else-if="entries.length === 0" class="text-secondary">
      {{ $t('article.editor.emptyList') }}
    </p>

    <template v-else>
      <table class="w-full table-auto">
        <thead class="border-b">
          <tr>
            <th class="py-2 text-left">{{ $t('article.editor.columnTitle') }}</th>
            <th class="py-2 text-left">{{ $t('article.editor.columnUri') }}</th>
            <th class="py-2 text-left">{{ $t('article.editor.columnStatus') }}</th>
            <th class="py-2 text-left">{{ $t('article.editor.columnUpdated') }}</th>
            <th class="py-2 text-left">{{ $t('article.editor.columnPublished') }}</th>
            <th class="py-2">{{ $t('common.actions') }}</th>
          </tr>
        </thead>
        <tbody class="border-b">
          <tr v-for="entry in entries" :key="entry.uri">
            <td class="py-2">{{ entry.title || entry.uri }}</td>
            <td class="py-2 text-secondary">/article/{{ entry.uri }}</td>
            <td class="py-2 text-secondary">{{ formatStatus(entry) }}</td>
            <td class="py-2 text-secondary">{{ formatUpdated(entry.updated_at) }}</td>
            <td class="py-2 text-secondary">{{ formatPublished(entry.published_at) }}</td>
            <td class="flex flex-wrap justify-center gap-2 py-2">
              <router-link
                class="rounded-xl border border-dotted px-4 py-px"
                :to="{ path: `/article/editor/${entry.uri}`, query: { page: $props.page } }"
              >
                {{ $t('article.editor.open') }}
              </router-link>
              <button
                class="rounded-xl border border-dotted px-4 py-px disabled:opacity-50"
                :disabled="busy || !entry.published_at"
                @click="unpublish(entry.uri)"
              >
                {{ $t('article.editor.unpublish') }}
              </button>
              <button
                class="rounded-xl border border-red-500/40 px-4 py-px text-red-400 disabled:opacity-50"
                :disabled="busy"
                @click="purge(entry.uri)"
              >
                {{ $t('common.delete') }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>

      <ArticlePager
        :page="$props.page"
        :page_size="$props.page_size"
        :total="total"
        :total_pages="total_pages"
        @change="$emit('page-change', $event)"
      ></ArticlePager>
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'

import ArticlePager from './ArticlePager.vue'
import { errorStatus, isArticle } from '@/helpers/article'
import { useToast } from '@/composables/useToast'
import {
  createArticleRaw,
  listArticlesAdminRaw,
  purgeArticleRaw,
  unpublishArticleRaw
} from '@/services/article'

import type { ArticleListEntry } from '@/interfaces/Article'

dayjs.extend(relativeTime)

interface Props {
  page?: number
  page_size?: number
}

const $props = withDefaults(defineProps<Props>(), {
  page: 1,
  page_size: 10
})

const $emit = defineEmits<{
  'page-change': [page: number]
}>()

const { t } = useI18n()
const router = useRouter()
const toast = useToast()

const entries = ref<ArticleListEntry[]>([])
const total = ref(0)
const total_pages = ref(0)
const loading = ref(true)
const busy = ref(false)
const conflict = ref(false)
const create_failed = ref(false)
const list_failed = ref(false)

async function refresh(): Promise<void> {
  loading.value = true
  list_failed.value = false

  try {
    const response = await listArticlesAdminRaw($props.page, $props.page_size).send()
    entries.value = Array.isArray(response?.data) ? response.data : []
    total.value = Number(response?.total) || 0
    total_pages.value = Number(response?.total_pages) || 0

    if (entries.value.length === 0 && $props.page > 1) {
      $emit('page-change', $props.page - 1)
    }
  } catch {
    entries.value = []
    total.value = 0
    total_pages.value = 0
    list_failed.value = true
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

function formatStatus(entry: ArticleListEntry): string {
  if (!entry.published_at) return t('article.editor.statusDraft')

  return dayjs(entry.published_at).isAfter(dayjs())
    ? t('article.editor.statusScheduled')
    : t('article.editor.statusPublished')
}

function formatUpdated(value: string): string {
  return dayjs(value).fromNow()
}

async function create(): Promise<void> {
  conflict.value = false
  create_failed.value = false

  if (!confirm(t('article.editor.confirmCreateBlank'))) return

  busy.value = true

  try {
    const response = await createArticleRaw({}).send()

    if (!isArticle(response) || !response.uri) {
      create_failed.value = true
      return
    }

    await router.push({ path: `/article/editor/${response.uri}`, query: { page: $props.page } })
  } catch (failure) {
    conflict.value = errorStatus(failure) === 409
    create_failed.value = !conflict.value
  } finally {
    busy.value = false
  }
}

async function unpublish(uri: string): Promise<void> {
  if (!confirm(t('article.editor.confirmUnpublish'))) return

  busy.value = true
  conflict.value = false
  create_failed.value = false

  try {
    await unpublishArticleRaw(uri).send()
    await refresh()
    toast.success(t('article.editor.unpublishedToast'))
  } catch {
    toast.error(t('article.editor.unpublishFailed'))
  } finally {
    busy.value = false
  }
}

async function purge(uri: string): Promise<void> {
  if (!confirm(t('article.editor.confirmDelete'))) return

  busy.value = true
  conflict.value = false
  create_failed.value = false

  try {
    await purgeArticleRaw(uri).send()
    await refresh()
    toast.success(t('article.editor.deleted'))
  } catch {
    toast.error(t('article.editor.deleteFailed'))
  } finally {
    busy.value = false
  }
}

watch(() => [$props.page, $props.page_size], refresh)

onMounted(refresh)
</script>
