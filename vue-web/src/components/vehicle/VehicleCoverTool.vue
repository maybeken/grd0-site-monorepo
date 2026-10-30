<template>
  <div class="flex flex-col gap-3">
    <p class="text-sm font-medium leading-relaxed text-foreground/60">
      {{ $t('grandTour.editor.cover.limitations') }}
    </p>

    <div
      v-if="status === 'unsupported'"
      class="rounded border border-red-500/40 bg-red-500/5 p-3 text-sm font-medium text-red-400"
      role="alert"
    >
      {{ $t('grandTour.editor.cover.unsupported') }}
    </div>

    <template v-else>
      <div class="flex flex-wrap items-center gap-3">
        <EditorButton variant="secondary" :disabled="working" @click="pick">
          {{ $t('grandTour.editor.cover.choosePhoto') }}
        </EditorButton>
        <EditorButton v-if="generated" variant="destructive" :disabled="working" @click="discard">
          {{ $t('grandTour.editor.cover.remove') }}
        </EditorButton>
        <input
          ref="file_input"
          type="file"
          accept="image/*"
          class="hidden"
          @change="onPhoto"
        />
      </div>

      <p
        v-if="working && progress !== null"
        class="text-sm font-medium text-foreground/60"
        aria-live="polite"
      >
        {{ $t('grandTour.editor.cover.modelProgress', { progress }) }}
      </p>
      <p v-else-if="working" class="text-sm font-medium text-foreground/60" aria-live="polite">
        {{ $t('grandTour.editor.cover.generating') }}
      </p>
      <div v-else-if="error_message" role="alert">
        <p class="text-sm font-medium text-red-400">{{ error_message }}</p>
        <p v-if="error_detail" class="mt-1 break-words text-xs text-foreground/60">
          {{ error_detail }}
        </p>
      </div>
      <p v-else-if="status === 'ready'" class="text-sm font-medium text-foreground/60">
        {{ $t('grandTour.editor.cover.modelReady') }}
      </p>

      <div v-if="generated" class="flex flex-col gap-3">
        <p class="text-lg font-medium">{{ $t('grandTour.editor.cover.variants') }}</p>
        <div class="grid grid-cols-1 gap-3 md:grid-cols-3">
          <figure class="flex flex-col gap-1">
            <div class="checker flex h-28 items-center justify-center rounded border border-accent/20">
              <img
                :src="generated.cutout_url"
                :alt="$t('grandTour.editor.cover.cutout')"
                class="h-full w-full object-contain"
              />
            </div>
            <figcaption class="text-sm font-medium text-foreground/60">
              {{ $t('grandTour.editor.cover.cutout') }}
            </figcaption>
          </figure>

          <figure class="flex flex-col gap-1">
            <div class="flex h-28 items-center justify-center rounded border border-accent/20 bg-background">
              <img
                :src="generated.silhouette_light_url"
                :alt="$t('grandTour.editor.cover.silhouetteLight')"
                class="h-full w-full object-contain"
              />
            </div>
            <figcaption class="text-sm font-medium text-foreground/60">
              {{ $t('grandTour.editor.cover.silhouetteLight') }}
            </figcaption>
          </figure>

          <figure class="flex flex-col gap-1">
            <div class="flex h-28 items-center justify-center rounded border border-accent/20 bg-foreground">
              <img
                :src="generated.silhouette_dark_url"
                :alt="$t('grandTour.editor.cover.silhouetteDark')"
                class="h-full w-full object-contain"
              />
            </div>
            <figcaption class="text-sm font-medium text-foreground/60">
              {{ $t('grandTour.editor.cover.silhouetteDark') }}
            </figcaption>
          </figure>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import EditorButton from '@/components/editor/ui/EditorButton.vue'
import { CUTOUT_MODEL_HOST, CUTOUT_MODEL_ID, CUTOUT_RUNTIME_HOST } from '@/helpers/vehicle'

import type { GeneratedCover } from '@/interfaces/Vehicle'

type WorkerMessage =
  | { type: 'progress'; progress: number }
  | { type: 'ready'; device: string }
  | { type: 'unsupported' }
  | {
      type: 'result'
      id: number
      cutout: Blob
      silhouette_light: Blob
      silhouette_dark: Blob
      width: number
      height: number
    }
  | { type: 'error'; id: number; message: string }

type ResultMessage = Extract<WorkerMessage, { type: 'result' }>

const emit = defineEmits<{
  generated: [cover: GeneratedCover]
  discarded: []
}>()

const { t } = useI18n()

const file_input = ref<HTMLInputElement | null>(null)
const generated = ref<GeneratedCover | null>(null)
const working = ref(false)
const progress = ref<number | null>(null)
const status = ref<'idle' | 'ready' | 'unsupported'>('idle')
const error_message = ref('')
const error_detail = ref('')

let worker: Worker | null = null
let next_id = 0
const pending = new Map<
  number,
  { resolve: (result: ResultMessage) => void; reject: (error: Error) => void }
>()

function ensure_worker(): Worker {
  if (worker) return worker

  const instance = new Worker(new URL('../../workers/cutout.worker.ts', import.meta.url), {
    type: 'module'
  })

  instance.onmessage = (event: MessageEvent<WorkerMessage>) => on_message(event.data)
  instance.onerror = (error) => {
    console.error('Cutout worker failed', error)
    working.value = false
    error_message.value = t('grandTour.editor.cover.generationFailed')
    error_detail.value = error.message || 'The cutout worker stopped unexpectedly.'
  }

  worker = instance

  return instance
}

function on_message(message: WorkerMessage): void {
  if (message.type === 'progress') {
    progress.value = message.progress
    return
  }

  if (message.type === 'ready') {
    status.value = 'ready'
    progress.value = null
    return
  }

  if (message.type === 'unsupported') {
    status.value = 'unsupported'
    working.value = false
    progress.value = null

    for (const entry of pending.values()) entry.reject(new Error('unsupported'))
    pending.clear()

    return
  }

  const entry = pending.get(message.id)
  if (!entry) return

  pending.delete(message.id)

  if (message.type === 'error') {
    entry.reject(new Error(message.message))
  } else {
    entry.resolve(message)
  }
}

function segment(blob: Blob): Promise<ResultMessage> {
  const id = (next_id += 1)
  const instance = ensure_worker()

  return new Promise<ResultMessage>((resolve, reject) => {
    pending.set(id, { resolve, reject })
    instance.postMessage({
      type: 'segment',
      id,
      blob,
      model_host: CUTOUT_MODEL_HOST,
      model_id: CUTOUT_MODEL_ID,
      runtime_host: CUTOUT_RUNTIME_HOST
    })
  })
}

function pick(): void {
  file_input.value?.click()
}

function revoke(): void {
  if (!generated.value) return

  URL.revokeObjectURL(generated.value.cutout_url)
  URL.revokeObjectURL(generated.value.silhouette_light_url)
  URL.revokeObjectURL(generated.value.silhouette_dark_url)
}

function apply(result: ResultMessage): void {
  revoke()

  const cover: GeneratedCover = {
    cutout: result.cutout,
    silhouette_light: result.silhouette_light,
    silhouette_dark: result.silhouette_dark,
    cutout_url: URL.createObjectURL(result.cutout),
    silhouette_light_url: URL.createObjectURL(result.silhouette_light),
    silhouette_dark_url: URL.createObjectURL(result.silhouette_dark)
  }

  generated.value = cover
  emit('generated', cover)
}

function discard(): void {
  revoke()
  generated.value = null
  emit('discarded')
}

async function onPhoto(event: Event): Promise<void> {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]

  input.value = ''

  if (!file || working.value) return

  working.value = true
  progress.value = null
  error_message.value = ''
  error_detail.value = ''

  try {
    apply(await segment(file))
  } catch (error) {
    console.error('Cutout generation failed', error)

    if (status.value !== 'unsupported') {
      error_message.value = t('grandTour.editor.cover.generationFailed')
      error_detail.value = error instanceof Error ? error.message : String(error)
    }
  } finally {
    working.value = false
    progress.value = null
  }
}

onMounted(() => {
  if (typeof WebAssembly === 'undefined') status.value = 'unsupported'
})

onBeforeUnmount(() => {
  revoke()
  worker?.terminate()
  worker = null
  pending.clear()
})
</script>

<style scoped>
.checker {
  background-image: conic-gradient(
    from 90deg at 50% 50%,
    var(--color-shade) 25%,
    var(--color-background) 0 50%,
    var(--color-shade) 0 75%,
    var(--color-background) 0
  );
  background-size: 16px 16px;
}
</style>
