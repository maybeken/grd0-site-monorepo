<template>
  <iframe ref="tldraw" class="w-full h-[80lvh] rounded-xl" src="/tldraw/index.html"></iframe>

  <div>
    <textarea class="w-full h-96 border rounded-xl px-2 text-xs" v-model="content"></textarea>
    <div class="flex py-2">
      <button class="px-4 border rounded-xl" @click="sendCommand('save')">Save</button>
      <button class="px-4 border rounded-xl" @click="sendCommand('load', { payload: base64Encode ? decompressFromBase64(content) : content })">Load</button>
      <div class="flex gap-2 px-4">
        <button class="min-w-8 min-h-8 bg-accent rounded-full" @click="base64Encode = !base64Encode">
          <Icon class="mx-auto my-auto" v-if="base64Encode" icon="mynaui:check" height="auto"></Icon>
        </button>
        <p class="my-auto text-sm">LZ Compressed + Base64</p>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { compressToBase64, decompressFromBase64 } from 'lz-string';
import { ref, onMounted, useTemplateRef } from 'vue';

const iframe = useTemplateRef('tldraw');
const content = ref();
const base64Encode = ref(true);

function sendCommand(command: string, payload: object = {}): void {
  iframe.value?.contentWindow?.postMessage({ event: "command", func: command, ...payload }, '*');
}

onMounted(() => {
  window.addEventListener(
    'message',
    (event) => {
      if (event.data.save) {
        content.value = base64Encode.value ? compressToBase64(JSON.stringify(event.data.save)) : JSON.stringify(event.data.save);
      }
    },
    false,
  );

  window.addEventListener(
    'keydown',
    (event) => {
      if ((event.ctrlKey || event.metaKey) && event.key === 's') {
        event.preventDefault();
        sendCommand('save');
      }
    },
    false,
  );
})
</script>
