<template>
  <div class="hidden lg:block fixed top-0 left-0 w-screen h-lvh z-50 bg-shade bg-opacity-50 backdrop-blur-xs motion-preset-fade-lg">
    <div class="flex w-full h-full" @click="$emit('close')">
      <div class="w-fit h-fit mx-auto my-auto p-4 rounded-2xl bg-shade text-secondary" @click="(event) => event.stopImmediatePropagation()">
        <div class="flex flex-row-reverse gap-4">
          <button @click="$emit('close')"><Icon icon="mynaui:x-circle" height="auto"></Icon></button>
        </div>
        <div class="h-screen">
          <img class="mx-auto h-[90%] object-scale-down rounded-xl" :src="`${ASSET_URL}${uri}`" />
        </div>
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
  close: [],
}>();
const { uri } = defineProps<Props>();
</script>