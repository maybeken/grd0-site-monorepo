<template>
  <div class="flex flex-col gap-2">
    <div class="relative" @mouseleave="interactive = false">
      <iframe
        ref="iframe"
        class="w-full h-[60vh] rounded-xl"
        :class="interactive ? '' : 'pointer-events-none'"
        src="/excalidraw/index.html?view=1"
        frameborder="0"
      ></iframe>

      <div
        v-if="!interactive && !failed"
        class="absolute inset-0 rounded-xl cursor-pointer"
        @click="activate"
      >
        <span
          class="absolute bottom-2 right-2 rounded-full border border-dotted bg-background/80 px-3 py-1 text-xs text-secondary"
        >
          {{ $t('article.block.drawingInteract') }}
        </span>
      </div>
    </div>

    <p v-if="loading" class="text-center text-sm text-secondary">
      {{ $t('article.block.drawingLoading') }}
    </p>
    <p v-else-if="failed" class="text-center text-sm text-red-400">
      {{ $t('article.block.drawingUnavailable') }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, useTemplateRef, watch } from 'vue'

import { articleSceneId } from '@/helpers/article'
import { getArticleDraftSceneRaw, getArticleSceneRaw } from '@/services/article'

import type { ArticleBlockDrawing } from '@/interfaces/Article'

interface Props {
  block: ArticleBlockDrawing
  uri: string
  draft?: boolean
}

const $props = defineProps<Props>()

const iframe = useTemplateRef('iframe')
const loading = ref(true)
const failed = ref(false)
const ready = ref(false)
const interactive = ref(false)

function activate(): void {
  interactive.value = true
  iframe.value?.contentWindow?.focus()
}

let load_started = false

async function loadScene(): Promise<void> {
  if (load_started) return
  load_started = true

  const scene_id = articleSceneId($props.block.scene_path)
  if (!scene_id) {
    failed.value = true
    loading.value = false
    return
  }

  try {
    const response = $props.draft
      ? await getArticleDraftSceneRaw($props.uri, scene_id).send()
      : await getArticleSceneRaw($props.uri, scene_id).send()

    const scene = response as { type?: string }
    if (!scene || typeof scene !== 'object' || scene.type !== 'excalidraw') {
      throw new Error('Invalid scene')
    }

    iframe.value?.contentWindow?.postMessage(
      { event: 'command', func: 'load', payload: JSON.stringify(scene) },
      window.location.origin
    )
  } catch {
    failed.value = true
  } finally {
    loading.value = false
  }
}

function onMessage(event: MessageEvent): void {
  if (event.source !== iframe.value?.contentWindow) return
  if (event.origin !== window.location.origin) return

  const data = event.data as { ready?: boolean; loadError?: string } | null
  if (!data || typeof data !== 'object') return

  if (data.ready === true) {
    ready.value = true
    loadScene()
  } else if (data.loadError) {
    failed.value = true
    loading.value = false
  }
}

onMounted(() => {
  window.addEventListener('message', onMessage)
})

onUnmounted(() => {
  window.removeEventListener('message', onMessage)
})

watch(
  () => $props.block.scene_path,
  (value) => {
    if (!value) return

    load_started = false
    failed.value = false
    loading.value = true

    if (ready.value) loadScene()
  }
)
</script>
