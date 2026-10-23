<template>
  <div class="flex min-h-lvh bg-background text-foreground">
    <EditorSidebar :open="sidebar_open" @close="sidebar_open = false"></EditorSidebar>

    <div class="flex min-w-0 flex-1 flex-col">
      <header class="flex items-center gap-3 border-b border-accent/20 px-4 py-2 lg:hidden">
        <button
          class="rounded border border-accent/40 p-1"
          :aria-label="$t('editor.openMenu')"
          @click="sidebar_open = true"
        >
          <Icon icon="mynaui:menu" />
        </button>
        <span class="text-sm font-semibold uppercase tracking-wide">{{ $t('editor.title') }}</span>
      </header>

      <main class="min-h-0 flex-1 overflow-y-auto p-4 md:p-6">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import EditorSidebar from './EditorSidebar.vue'

const route = useRoute()
const sidebar_open = ref(false)

watch(
  () => route.fullPath,
  () => {
    sidebar_open.value = false
  }
)
</script>
