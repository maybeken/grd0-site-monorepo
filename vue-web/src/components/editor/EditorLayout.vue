<template>
  <div class="flex min-h-lvh bg-background text-foreground">
    <EditorSidebar :open="sidebar_open" @close="sidebar_open = false"></EditorSidebar>

    <div class="flex min-w-0 flex-1 flex-col">
      <header class="flex items-center gap-3 border-b border-accent/20 px-4 py-2 lg:hidden">
        <EditorButton
          variant="icon"
          :label="$t('editor.openMenu')"
          @click="sidebar_open = true"
        >
          <Icon icon="mynaui:menu"></Icon>
        </EditorButton>
        <span class="text-sm font-medium uppercase tracking-wide text-foreground/60">
          {{ $t('editor.title') }}
        </span>
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
import EditorButton from './ui/EditorButton.vue'
import Icon from '@/components/shared/Icon.vue'

const route = useRoute()
const sidebar_open = ref(false)

watch(
  () => route.fullPath,
  () => {
    sidebar_open.value = false
  }
)
</script>
