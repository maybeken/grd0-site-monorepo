<template>
  <ArticleEditor v-if="uri" :uri="uri" @saved="onSaved"></ArticleEditor>

  <EditorPage v-else>
    <EditorError :message="$t('editor.notFound')">
      <template #actions>
        <EditorButton variant="secondary" to="/editor/articles">
          {{ $t('editor.ui.backTo', { section: $t('editor.articles') }) }}
        </EditorButton>
      </template>
    </EditorError>
  </EditorPage>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import ArticleEditor from '@/components/editor/article/ArticleEditor.vue'
import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorError from '@/components/editor/ui/EditorError.vue'
import EditorPage from '@/components/editor/ui/EditorPage.vue'

const route = useRoute()
const router = useRouter()

const uri = computed(() => {
  const value = route.params.uri

  if (typeof value === 'string') return value
  if (Array.isArray(value)) return value[0] || ''

  return ''
})

function onSaved(value: string): void {
  if (value === uri.value) return

  router.replace({ path: `/editor/articles/${value}`, query: route.query })
}
</script>
