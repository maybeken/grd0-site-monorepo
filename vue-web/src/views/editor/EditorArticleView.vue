<template>
  <div class="flex flex-col gap-4">
    <ArticleEditor v-if="uri" :uri="uri" @saved="onSaved"></ArticleEditor>

    <div v-else class="flex flex-col items-start gap-3">
      <p class="text-red-400">{{ $t('editor.notFound') }}</p>
      <router-link to="/editor/articles" class="rounded-full border border-dotted px-3 py-1">
        {{ $t('editor.backToList') }}
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import ArticleEditor from '@/components/editor/article/ArticleEditor.vue'

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
