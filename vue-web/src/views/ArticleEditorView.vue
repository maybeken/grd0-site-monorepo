<template>
  <div class="flex flex-col gap-8">
    <p class="text-2xl font-bold">{{ $t('article.editor.title') }}<CursorBlink /></p>

    <div v-if="authorized" class="flex flex-col gap-8">
      <ArticleList
        ref="list"
        :selected="selected"
        @open="selected = $event"
        @unpublished="onUnpublished"
      ></ArticleList>

      <ArticleEditor
        v-if="selected"
        :key="selected"
        :initial-uri="selected"
        @saved="onSaved"
      ></ArticleEditor>
      <p v-else class="text-secondary">{{ $t('article.editor.selectArticle') }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, useTemplateRef } from 'vue'

import ArticleList from '@/components/editor/article/ArticleList.vue'
import ArticleEditor from '@/components/editor/article/ArticleEditor.vue'
import { isLocalEnvironment } from '@/helpers/environment'

const authorized = ref(false)
const selected = ref('')
const list = useTemplateRef<InstanceType<typeof ArticleList>>('list')

onMounted(() => {
  if (!isLocalEnvironment() && !sessionStorage.getItem('jwt_token')) {
    window.location.href = '/redirect/auth/login'
    return
  }

  authorized.value = true
})

function onSaved(uri: string): void {
  if (uri !== selected.value) selected.value = uri

  list.value?.refresh()
}

function onUnpublished(uri: string): void {
  list.value?.refresh()
}
</script>
