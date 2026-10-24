<template>
  <section class="rounded-xl border border-accent/20 bg-background p-4">
    <div
      v-if="title || description || $slots.actions"
      class="mb-4 flex flex-wrap items-start justify-between gap-2"
    >
      <div class="flex min-w-0 flex-col gap-1">
        <h2 v-if="title" class="editor-section-title font-bold">{{ title }}</h2>
        <p v-if="description" class="text-sm font-medium text-foreground/60">{{ description }}</p>
      </div>
      <div v-if="$slots.actions" class="flex flex-wrap items-center gap-2">
        <slot name="actions"></slot>
      </div>
    </div>

    <slot></slot>
  </section>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    title?: string
    description?: string
  }>(),
  {
    title: '',
    description: ''
  }
)
</script>

<style scoped>
.editor-section-title {
  font-size: 1.25rem;
  line-height: 1.3;
}
</style>
