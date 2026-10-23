<template>
  <div
    class="flex items-center justify-between gap-4 rounded-lg border border-accent/20 bg-background p-3"
  >
    <div class="flex min-w-0 flex-col">
      <span class="truncate text-sm font-semibold">{{ title }}</span>
      <slot></slot>
    </div>
    <div class="flex shrink-0 gap-2">
      <RouterLink
        :to="to"
        class="rounded border border-accent/40 px-2 py-1 text-xs hover:bg-accent/10"
      >
        {{ $t('common.edit') }}
      </RouterLink>
      <button
        class="rounded border border-red-500/40 px-2 py-1 text-xs text-red-400 hover:bg-red-500/10"
        @click="$emit('delete')"
      >
        {{ $t('common.delete') }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  to: string
}>()

defineEmits<{
  delete: []
}>()
</script>
