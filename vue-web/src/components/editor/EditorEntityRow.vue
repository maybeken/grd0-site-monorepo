<template>
  <div
    class="flex items-center justify-between gap-4 rounded-lg border border-accent/20 bg-background p-3"
  >
    <div class="flex min-w-0 items-center gap-3">
      <div v-if="$slots.leading" class="shrink-0">
        <slot name="leading"></slot>
      </div>
      <div class="flex min-w-0 flex-col">
        <span class="truncate text-lg font-medium" :title="title">{{ title }}</span>
        <slot></slot>
      </div>
    </div>
    <div class="flex shrink-0 flex-wrap gap-2">
      <EditorButton :to="to" size="sm" variant="secondary">{{ $t('common.edit') }}</EditorButton>
      <EditorButton size="sm" variant="destructive" @click="$emit('delete')">
        {{ $t('common.delete') }}
      </EditorButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import EditorButton from './ui/EditorButton.vue'

defineProps<{
  title: string
  to: string
}>()

defineEmits<{
  delete: []
}>()
</script>
