<template>
  <div class="flex flex-col gap-8">
    <p class="text-2xl font-bold">{{ $t('article.editor.title') }}<CursorBlink /></p>

    <div v-if="authorized" class="flex flex-col gap-8">
      <ArticleList
        ref="list"
        :selected="selected"
        @open="openArticle"
        @unpublished="onUnpublished"
        @deleted="onDeleted"
      ></ArticleList>

      <ArticleEditor
        v-if="selected"
        :key="selected"
        :initial-uri="selected"
        @saved="onSaved"
        @dirty-change="editor_dirty = $event"
      ></ArticleEditor>
      <p v-else class="text-secondary">{{ $t('article.editor.selectArticle') }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'

import ArticleList from '@/components/editor/article/ArticleList.vue'
import ArticleEditor from '@/components/editor/article/ArticleEditor.vue'
import { isLocalEnvironment } from '@/helpers/environment'

const { t } = useI18n()

const authorized = ref(false)
const editor_dirty = ref(false)
const selected = ref('')
const list = useTemplateRef<InstanceType<typeof ArticleList>>('list')

onMounted(() => {
  if (!isLocalEnvironment() && !sessionStorage.getItem('jwt_token')) {
    window.location.href = '/redirect/auth/login'
    return
  }

  authorized.value = true
})

function openArticle(uri: string): void {
  if (uri === selected.value) return
  if (editor_dirty.value && !window.confirm(t('article.editor.confirmDiscard'))) return

  editor_dirty.value = false
  selected.value = uri
}

function onSaved(uri: string): void {
  if (uri !== selected.value) selected.value = uri

  list.value?.refresh()
}

function onUnpublished(uri: string): void {
  list.value?.refresh()
}

function onDeleted(uri: string): void {
  if (selected.value !== uri) return

  selected.value = ''
  editor_dirty.value = false
}
</script>
