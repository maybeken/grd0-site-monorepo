<template>
  <EditorPage>
    <EditorPageHeader
      :title="$t('editor.articles')"
      :description="$t('article.editor.listHint')"
    >
      <template #actions>
        <EditorButton variant="primary" @click="article_list?.create()">
          <Icon icon="mynaui:plus" height="1rem"></Icon>
          {{ $t('article.editor.newArticle') }}
        </EditorButton>
      </template>
    </EditorPageHeader>

    <ArticleList
      ref="article_list"
      :page="page"
      :page_size="PAGE_SIZE"
      @page-change="setPage"
    ></ArticleList>
  </EditorPage>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import ArticleList from '@/components/editor/article/ArticleList.vue'
import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorPage from '@/components/editor/ui/EditorPage.vue'
import EditorPageHeader from '@/components/editor/ui/EditorPageHeader.vue'
import Icon from '@/components/shared/Icon.vue'

const PAGE_SIZE = 10

const route = useRoute()
const router = useRouter()
const article_list = ref<InstanceType<typeof ArticleList> | null>(null)

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
