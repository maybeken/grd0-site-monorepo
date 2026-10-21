<template>
  <div v-if="authorized" class="flex flex-col gap-8">
    <header class="flex flex-col gap-2">
      <h1 class="text-2xl font-bold">{{ $t('article.editor.title') }}<CursorBlink /></h1>
      <p class="text-secondary">{{ $t('article.editor.listHint') }}</p>
    </header>

    <ArticleList :page="page" :page_size="PAGE_SIZE" @page-change="setPage"></ArticleList>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import ArticleList from '@/components/editor/article/ArticleList.vue'
import { isLocalEnvironment } from '@/helpers/environment'

const PAGE_SIZE = 10

const route = useRoute()
const router = useRouter()
const authorized = ref(false)

const page = computed(() => {
  const value = Number(route.query.page)

  return Number.isInteger(value) && value > 0 ? value : 1
})

function setPage(value: number): void {
  if (value === page.value) return

  router.push({
    path: '/article/editor',
    query: { ...route.query, page: String(value) }
  })
}

onMounted(() => {
  if (!isLocalEnvironment() && !sessionStorage.getItem('jwt_token')) {
    window.location.href = '/redirect/auth/login'
    return
  }

  authorized.value = true
})
</script>
