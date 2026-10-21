<template>
  <div v-if="authorized" class="flex flex-col gap-4">
    <ArticleEditor v-if="uri" :uri="uri" @saved="onSaved"></ArticleEditor>

    <div v-else class="flex flex-col items-start gap-3">
      <p class="text-red-400">{{ $t('article.editor.loadFailed') }}</p>
      <router-link to="/article/editor" class="rounded-full border border-dotted px-3 py-1">
        {{ $t('article.editor.backToList') }}
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import ArticleEditor from '@/components/editor/article/ArticleEditor.vue'
import { isLocalEnvironment } from '@/helpers/environment'

const route = useRoute()
const router = useRouter()
const authorized = ref(false)

const uri = computed(() => {
  const value = route.params.uri

  if (typeof value === 'string') return value
  if (Array.isArray(value)) return value[0] || ''

  return ''
})

function onSaved(value: string): void {
  if (value === uri.value) return

  router.replace({ path: `/article/editor/${value}`, query: route.query })
}

onMounted(() => {
  if (!isLocalEnvironment() && !sessionStorage.getItem('jwt_token')) {
    window.location.href = '/redirect/auth/login'
    return
  }

  authorized.value = true
})
</script>
