<template>
  <div v-if="article" class="flex flex-col gap-6 pb-8">
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
import { useRoute } from 'vue-router'

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
</script>
