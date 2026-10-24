<template>
  <div class="flex flex-col gap-2">
    <p class="text-lg font-medium text-foreground/60">{{ $t('article.editor.drawingHint') }}</p>

    <div class="relative" @mouseleave="interactive = false">
      <iframe
        ref="iframe"
        class="h-[60vh] w-full rounded-xl"
        :class="interactive ? '' : 'pointer-events-none'"
        src="/excalidraw/index.html"
        frameborder="0"
      ></iframe>

      <div
        v-if="!interactive && !load_failed"
        class="absolute inset-0 cursor-pointer rounded-xl"
        @click="activate"
      >
        <span
          class="absolute bottom-2 right-2 rounded-full border border-dotted bg-background/80 px-3 py-1 text-sm font-medium text-foreground/60"
        >
          {{ $t('article.block.drawingInteract') }}
        </span>
      </div>

      <p
        v-if="loading"
        class="pointer-events-none absolute inset-0 grid place-content-center text-lg font-medium text-foreground/60"
      >
        {{ $t('article.block.drawingLoading') }}
      </p>
    </div>

    <div class="flex flex-wrap items-center gap-2">
      <EditorButton
        variant="secondary"
        size="sm"
        :disabled="!ready || saving"
        @click="saveScene"
      >
        {{ $t('article.editor.saveDrawing') }}
      </EditorButton>
      <span v-if="saving" class="text-sm font-medium text-foreground/60">
        {{ $t('common.saving') }}
      </span>
      <span v-if="save_failed" class="text-sm font-medium text-red-400">
        {{ $t('article.editor.sceneFailed') }}
      </span>
      <span v-if="load_failed" class="text-sm font-medium text-red-400">
        {{ $t('article.editor.sceneLoadFailed') }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, useTemplateRef } from 'vue'

import EditorButton from '@/components/editor/ui/EditorButton.vue'
import { articleSceneId } from '@/helpers/article'
import { getArticleDraftSceneRaw, putArticleSceneRaw } from '@/services/article'

import type { ArticleBlockDrawing } from '@/interfaces/Article'

interface Props {
  block: ArticleBlockDrawing
  uri: string
}

const $props = defineProps<Props>()
const $emit = defineEmits<{
  requestSave: []
}>()

const iframe = useTemplateRef('iframe')
const ready = ref(false)
const loading = ref(true)
const saving = ref(false)
const save_failed = ref(false)
const load_failed = ref(false)
const interactive = ref(false)

function activate(): void {
  interactive.value = true
  iframe.value?.contentWindow?.focus()
}

const scene_id = ref(articleSceneId($props.block.scene_path) || crypto.randomUUID())

let save_promise: Promise<string | null> | null = null
let last_saved_at = 0
let last_saved_path: string | null = null
let pending_resolve: ((scene: object) => void) | null = null
let pending_reject: ((error: Error) => void) | null = null
let pending_timer: number | undefined
let scene_loaded = false

function post(message: unknown): void {
  iframe.value?.contentWindow?.postMessage(message, window.location.origin)
}

async function loadScene(): Promise<void> {
  if (scene_loaded) return
  scene_loaded = true

  const scene_id = articleSceneId($props.block.scene_path)
  if (!scene_id) return

  try {
    const response = await getArticleDraftSceneRaw($props.uri, scene_id).send()
    const scene = response as { type?: string }
    if (!scene || typeof scene !== 'object' || scene.type !== 'excalidraw') {
      throw new Error('Invalid scene')
    }

    post({ event: 'command', func: 'load', payload: JSON.stringify(scene) })
  } catch {
    load_failed.value = true
  }
}

function requestScene(): Promise<object> {
  if (!ready.value) return Promise.reject(new Error('Drawing is not ready'))

  return new Promise((resolve, reject) => {
    pending_resolve = resolve
    pending_reject = reject
    pending_timer = window.setTimeout(() => {
      pending_resolve = null
      pending_reject = null
      reject(new Error('Drawing save timed out'))
    }, 5000)

    post({ event: 'command', func: 'save' })
  })
}

async function doSaveScene(): Promise<string | null> {
  if (load_failed.value) {
    save_failed.value = true
    return null
  }

  saving.value = true
  save_failed.value = false

  try {
    const scene = await requestScene()
    const response = await putArticleSceneRaw($props.uri, scene_id.value, scene).send()

    if (!response?.scene_path) throw new Error('Missing scene path')

    $props.block.scene_path = response.scene_path
    return response.scene_path
  } catch {
    save_failed.value = true
    return null
  } finally {
    saving.value = false
  }
}

function saveScene(): Promise<string | null> {
  if (save_promise) return save_promise
  if (last_saved_path && Date.now() - last_saved_at < 500) {
    return Promise.resolve(last_saved_path)
  }

  save_promise = doSaveScene()
    .then((scene_path) => {
      if (scene_path) {
        last_saved_path = scene_path
        last_saved_at = Date.now()
      }
      return scene_path
    })
    .finally(() => {
      save_promise = null
    })

  return save_promise
}

function onMessage(event: MessageEvent): void {
  if (event.source !== iframe.value?.contentWindow) return
  if (event.origin !== window.location.origin) return

  const data = event.data as
    | { ready?: boolean; save?: unknown; loadError?: string; requestSave?: boolean }
    | null
  if (!data || typeof data !== 'object') return

  if (data.requestSave === true) {
    $emit('requestSave')
    return
  }

  if (data.ready === true) {
    ready.value = true
    loading.value = false
    loadScene()
    return
  }

  if (typeof data.save !== 'undefined') {
    window.clearTimeout(pending_timer)

    const resolve = pending_resolve
    pending_resolve = null
    pending_reject = null

    if (resolve && data.save && typeof data.save === 'object') {
      resolve(data.save as object)
    }
    return
  }

  if (data.loadError) {
    load_failed.value = true
  }
}

onMounted(() => {
  window.addEventListener('message', onMessage)
})

onUnmounted(() => {
  window.removeEventListener('message', onMessage)
  window.clearTimeout(pending_timer)
  pending_reject?.(new Error('Drawing editor closed'))
  pending_resolve = null
  pending_reject = null
})

defineExpose({ saveScene })
</script>
