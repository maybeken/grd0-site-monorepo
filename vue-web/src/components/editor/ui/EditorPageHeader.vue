<template>
  <header class="flex flex-wrap items-start justify-between gap-4">
    <div class="flex min-w-0 flex-col gap-2">
      <RouterLink
        v-if="backTo"
        :to="backTo"
        class="inline-flex w-fit items-center gap-1 text-lg font-medium text-foreground/60 hover:text-foreground"
      >
        <Icon icon="mynaui:arrow-left" height="1rem"></Icon>
        {{ backLabel }}
      </RouterLink>

      <h1 v-if="!titleSrOnly" :id="titleId" class="editor-page-title font-bold">
        {{ title }}<CursorBlink />
      </h1>
      <h1 v-else :id="titleId" class="sr-only">{{ title }}</h1>

      <p v-if="description" class="text-lg font-medium text-foreground/60">{{ description }}</p>
    </div>

    <div v-if="$slots.actions" class="flex flex-wrap items-center gap-2">
      <slot name="actions"></slot>
    </div>
  </header>
</template>

<script setup lang="ts">
import CursorBlink from '@/components/shared/CursorBlink.vue'
import Icon from '@/components/shared/Icon.vue'

import type { RouteLocationRaw } from 'vue-router'

withDefaults(
  defineProps<{
    title: string
    description?: string
    backTo?: RouteLocationRaw
    backLabel?: string
    titleId?: string
    titleSrOnly?: boolean
  }>(),
  {
    description: '',
    backTo: '',
    backLabel: '',
    titleId: 'editor-page-title',
    titleSrOnly: false
  }
)
</script>

<style scoped>
.editor-page-title {
  font-size: 1.5rem;
  line-height: 1.2;
}

@media (min-width: 768px) {
  .editor-page-title {
    font-size: 1.875rem;
  }
}
</style>
