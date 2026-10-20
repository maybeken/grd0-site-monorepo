<template>
  <header class="flex flex-col gap-2 py-4">
    <div class="flex flex-wrap items-center gap-3">
      <h1 class="text-4xl font-black">{{ $props.article.title || $props.article.uri }}</h1>
      <span
        v-if="$props.draft"
        class="rounded-full border border-accent px-3 py-1 text-sm text-secondary"
      >
        {{ $t('article.header.draft') }}
      </span>
    </div>

    <p v-if="$props.article.subtitle" class="text-xl text-secondary">
      {{ $props.article.subtitle }}
    </p>

    <div class="flex flex-wrap gap-x-4 gap-y-1 text-sm text-secondary">
      <span>
        {{
          $t('article.header.by', {
            author: $props.article.author?.display_name || $props.article.author?.email
          })
        }}
      </span>
      <span v-if="$props.article.published_at">
        {{
          $t('article.header.publishedAt', {
            date: formatDate($props.article.published_at)
          })
        }}
      </span>
      <span v-if="$props.article.updated_at">
        {{
          $t('article.header.updatedAt', {
            date: formatDate($props.article.updated_at)
          })
        }}
      </span>
    </div>
  </header>
</template>

<script setup lang="ts">
import dayjs from 'dayjs'
import LocalizedFormat from 'dayjs/plugin/localizedFormat'

import type { Article } from '@/interfaces/Article'

dayjs.extend(LocalizedFormat)

interface Props {
  article: Article
  draft?: boolean
}

const $props = defineProps<Props>()

function formatDate(value: string): string {
  return dayjs(value).format('LLL')
}
</script>
