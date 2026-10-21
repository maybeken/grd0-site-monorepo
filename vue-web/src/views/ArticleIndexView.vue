<template>
  <div class="flex flex-col gap-8">
    <header class="flex flex-col gap-2">
      <h1 class="text-2xl font-bold">{{ $t('article.index.title') }}<CursorBlink /></h1>
    </header>

    <template v-if="loading">
      <ArticleSummary v-for="i in 5" :key="i" :loading="loading"></ArticleSummary>
    </template>

    <p v-else-if="article_entries.length === 0" class="text-secondary">
      {{ $t('article.index.empty') }}
    </p>

    <ArticleSummary
      v-for="entry in article_entries"
      :key="entry.uri"
      :loading="loading"
      :entry="entry"
    ></ArticleSummary>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import ArticleSummary from '@/components/article/ArticleSummary.vue'
import { listArticles } from '@/services/article'

const response = listArticles()
const articles = response.data
const loading = response.loading

const article_entries = computed(() => (Array.isArray(articles.value) ? articles.value : []))
</script>
