<template>
  <iframe
    ref="excalidraw"
    class="w-full h-[80lvh] rounded-xl"
    src="/excalidraw/index.html"
  ></iframe>

  <div>
    <textarea class="w-full h-96 border rounded-xl px-2 text-xs" v-model="content"></textarea>
    <div class="flex gap-2 py-2">
      <button class="px-4 border rounded-xl" @click="sendCommand('save')">
        {{ $t('common.save') }}
      </button>
      <button
        class="px-4 border rounded-xl"
        :disabled="!content"
        @click="
          sendCommand('load', { payload: base64Encode ? decompressFromBase64(content) : content })
        "
      >
        {{ $t('common.load') }}
      </button>
      <div class="flex gap-2">
        <button
          class="min-w-8 min-h-8 bg-accent rounded-full"
          :class="base64Encode ? 'bg-green-800' : 'bg-accent'"
          @click="base64Encode = !base64Encode"
        >
          <Icon
            class="mx-auto my-auto"
            v-if="base64Encode"
            icon="mynaui:check"
            height="1.5rem"
          ></Icon>
        </button>
        <p class="my-auto text-sm">{{ $t('tool.excalidraw.lzCompressed') }}</p>
      </div>
    </div>
    <p v-if="loadError" class="text-red-400 text-sm py-1">
      {{ $t('tool.excalidraw.loadError') }}
    </p>
  </div>
</template>

<script lang="ts" setup>
import { compressToBase64, decompressFromBase64 } from 'lz-string'
import { onMounted, onUnmounted, ref, useTemplateRef } from 'vue'

const iframe = useTemplateRef('excalidraw')
const content = ref('')
const base64Encode = ref(true)
const loadError = ref(false)

function sendCommand(command: string, payload: object = {}): void {
  loadError.value = false
  iframe.value?.contentWindow?.postMessage({ event: 'command', func: command, ...payload }, '*')
}

function onMessage(event: MessageEvent): void {
  if (event.source !== iframe.value?.contentWindow) return

  if (event.data?.loadError) {
    loadError.value = true
  }

  if (typeof event.data?.save !== 'undefined') {
    content.value = base64Encode.value
      ? compressToBase64(JSON.stringify(event.data.save))
      : JSON.stringify(event.data.save)
  }
}

function onKeyDown(event: KeyboardEvent): void {
  if ((event.ctrlKey || event.metaKey) && event.key === 's') {
    event.preventDefault()
    sendCommand('save')
  }
}

onMounted(() => {
  window.addEventListener('message', onMessage, false)
  window.addEventListener('keydown', onKeyDown, false)
})

onUnmounted(() => {
  window.removeEventListener('message', onMessage, false)
  window.removeEventListener('keydown', onKeyDown, false)
})
</script>
