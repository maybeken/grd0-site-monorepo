<template>
  <div class="flex flex-col gap-8">
    <header class="flex flex-col gap-2">
      <h1 class="text-2xl font-bold">{{ $t('article.editor.title') }}<CursorBlink /></h1>
      <p class="text-secondary">{{ $t('article.editor.listHint') }}</p>
    </header>

    <ArticleList :page="page" :page_size="PAGE_SIZE" @page-change="setPage"></ArticleList>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import ArticleList from '@/components/editor/article/ArticleList.vue'

const PAGE_SIZE = 10

const route = useRoute()
const router = useRouter()

const page = computed(() => {
  const value = Number(route.query.page)

  return Number.isInteger(value) && value > 0 ? value : 1
})

function setPage(value: number): void {
  if (value === page.value) return

  router.push({
    path: '/editor/articles',
    query: { ...route.query, page: String(value) }
  })
}
</script>
