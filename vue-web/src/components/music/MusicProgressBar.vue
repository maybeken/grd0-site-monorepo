<template>
  <div class="relative w-full h-2 cursor-pointer">
    <div class="absolute left-0 h-2 rounded-2xl w-full z-20" @click.stopPropagation="seek"></div>
    <div
      class="absolute left-0 z-10 h-2 rounded-2xl bg-linear-to-r from-indigo-500 from-0% via-sky-500 via-50% to-emerald-500 to-100% text-right"
      :style="`width: ${(progress / duration) * 100}%`"
    ></div>
    <div class="absolute left-0 h-2 rounded-2xl bg-accent w-full"></div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  progress: number
  duration: number
}

const { progress = 0, duration = 1 } = defineProps<Props>()

const $emit = defineEmits(['seek'])

function seek(e: any) {
  $emit('seek', e.layerX / e.target.clientWidth)
}
</script>
