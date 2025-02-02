<template>
  <iframe ref="tldraw" class="w-full h-[80lvh] rounded-xl" src="/tldraw/index.html"></iframe>
  <textarea class="w-full h-96 border rounded-xl px-2 text-xs" v-model="content"></textarea>
  <div class="flex">
    <button class="px-4 border rounded-xl" @click="sendCommand('save')">Save</button>
    <button class="px-4 border rounded-xl" @click="sendCommand('load', { payload: content })">Load</button>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, useTemplateRef } from 'vue';

const iframe = useTemplateRef('tldraw');
const content = ref();

function sendCommand(command: string, payload: object = {}): void {
  iframe.value?.contentWindow?.postMessage({ event: "command", func: command, ...payload }, '*');
}

onMounted(() => {
  window.addEventListener(
    'message',
    (event) => {
      if (event.data.save) {
        content.value = JSON.stringify(event.data.save);
      }
    },
    false,
  );
})
</script>
