<template>
  <div class="relative h-12">
    <div class="absolute top-0 h-12 z-10 w-full md:w-96 max-w-full">
      <div
        class="flex px-4 py-2 w-full rounded-lg font-black bg-background border border-dotted border-foreground cursor-default z-20"
        @click="expanded = !expanded"
      >
        <span class="grow">{{ selected_title }}</span>
        <Icon icon="mynaui:chevron-down" height="auto" />
      </div>
      <div
        v-for="(option, idx) in options"
        :key="option.id"
        class="px-4 py-2 w-full bg-background bg-opacity-50 backdrop-blur-xs border border-dotted border-foreground cursor-pointer motion-preset-slide-down motion-duration-500"
        :class="optionClass(option.id, idx)"
        @click="pick(option.id)"
      >
        {{ option.title }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

import type { Playlist } from '@/interfaces/Music'

interface Props {
  playlists: Playlist[]
  selected: string
}

const $props = defineProps<Props>()
const $emit = defineEmits<{
  select: [value: string]
}>()

const expanded = ref(false)

const options = computed(() =>
  $props.playlists.map((playlist) => ({
    id: playlist.id ?? '',
    title: playlist.title
  }))
)

const selected_title = computed(
  () => options.value.find((option) => option.id === $props.selected)?.title ?? ''
)

function optionClass(id: string, idx: number): string[] {
  if (!expanded.value) return ['hidden']

  const style: string[] = []

  if (idx === 0) style.push('rounded-t-lg')
  if (idx === options.value.length - 1) style.push('rounded-b-lg')
  if (id === $props.selected) style.push('font-black')

  return style
}

function pick(id: string): void {
  $emit('select', id)
  expanded.value = false
}
</script>
