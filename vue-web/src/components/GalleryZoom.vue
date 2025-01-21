<template>
  <div class="hidden lg:block fixed top-0 left-0 w-screen h-screen z-50 bg-shade bg-opacity-50 backdrop-blur-sm motion-preset-fade-lg">
    <div class="flex w-full h-full" @click="$emit('close')">
      <div class="w-2/3 h-fit mx-auto my-auto p-4 rounded-2xl bg-shade text-secondary" @click="(event) => event.stopImmediatePropagation()">
        <div class="flex flex-row-reverse gap-4">
          <button @click="$emit('close')"><Icon icon="mynaui:x-circle" height="auto"></Icon></button>
        </div>
        <img class="w-full rounded-xl" :src="`${ASSET_URL}${uri}`" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const ASSET_URL = import.meta.env.VITE_ASSETS_URL;

interface Props {
  uri: string,
}

const $emit = defineEmits<{
  close: null,
}>();
const { uri } = defineProps<Props>();
</script>