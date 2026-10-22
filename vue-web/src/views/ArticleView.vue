<template>
  <div v-if="article" class="flex flex-col gap-6 pb-8 px-2 md:px-0">
    <p
      v-if="not_published"
      class="rounded-xl border border-yellow-400 bg-yellow-200 py-3 text-black text-center"
    >
      {{ $t('article.header.unpublishedNotice') }}
    </p>
    <ArticleHeader :article="article"></ArticleHeader>
    <BlockRenderer :blocks="article.blocks" :uri="article.uri"></BlockRenderer>
  </div>

  <div v-else-if="!loading" class="flex flex-col gap-4">
    <NotFound></NotFound>
  </div>

  <div v-else class="flex flex-col gap-4">
    <Skeleton h="xl" w="2/3" :loading="true"></Skeleton>
    <Skeleton h="md" w="full" :loading="true"></Skeleton>
    <Skeleton h="lg" w="full" :loading="true"></Skeleton>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import dayjs from 'dayjs'

import { getArticle } from '@/services/article'

const $route = useRoute()
const slug = $route.params.slug

let uri = ''
if (typeof slug === 'string') {
  uri = slug
} else if (Array.isArray(slug) && slug.length > 0) {
  uri = slug[0]!
}

const response = getArticle(uri)
const article = response?.data
const loading = response?.loading

const not_published = computed(() => {
  const value = article?.value?.published_at
  return !value || dayjs(value).isAfter(dayjs())
})
</script>
