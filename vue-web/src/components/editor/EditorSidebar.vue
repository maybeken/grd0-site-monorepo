<template>
  <template v-if="open">
    <div class="fixed inset-0 z-40 bg-black/60 lg:hidden" @click="$emit('close')"></div>
  </template>

  <aside
    class="fixed inset-y-0 left-0 z-50 flex w-64 shrink-0 flex-col border-r border-accent/20 bg-background transition-transform lg:static lg:z-auto lg:translate-x-0"
    :class="open ? 'translate-x-0' : '-translate-x-full'"
  >
    <div class="flex items-center justify-between px-4 py-4">
      <RouterLink to="/" class="text-lg font-extrabold uppercase tracking-wide">
        grd0.net
      </RouterLink>
      <button
        class="rounded border border-accent/40 p-1 lg:hidden"
        :aria-label="$t('editor.closeMenu')"
        @click="$emit('close')"
      >
        <Icon icon="mynaui:x" />
      </button>
    </div>

    <nav class="flex flex-1 flex-col gap-6 overflow-y-auto px-3 pb-6">
      <div v-for="group in groups" :key="group.label">
        <p class="px-2 pb-2 text-xs font-semibold uppercase tracking-widest text-secondary">
          {{ $t(group.label) }}
        </p>
        <RouterLink
          v-for="item in group.items"
          :key="item.to"
          :to="item.to"
          class="block rounded px-2 py-2 text-sm transition-colors"
          :class="
            isActive(item)
              ? 'bg-accent text-background font-semibold'
              : 'text-foreground/70 hover:bg-accent/10 hover:text-foreground'
          "
          :aria-current="isActive(item) ? 'page' : undefined"
        >
          {{ $t(item.label) }}
        </RouterLink>
      </div>
    </nav>

    <div class="border-t border-accent/20 p-3">
      <RouterLink to="/" class="block rounded px-2 py-2 text-sm text-secondary hover:bg-accent/10">
        {{ $t('editor.viewSite') }}
      </RouterLink>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'

defineProps<{
  open: boolean
}>()

defineEmits<{
  close: []
}>()

const route = useRoute()

const groups = [
  {
    label: 'editor.content',
    items: [
      { label: 'editor.articles', to: '/editor/articles', match: '/editor/articles' },
      { label: 'editor.gallery', to: '/editor/gallery', match: '/editor/gallery' }
    ]
  },
  {
    label: 'editor.coffee',
    items: [
      { label: 'editor.beans', to: '/editor/coffee/beans', match: '/editor/coffee/beans' },
      {
        label: 'editor.equipment',
        to: '/editor/coffee/equipment',
        match: '/editor/coffee/equipment'
      },
      { label: 'editor.tastings', to: '/editor/coffee/tastings', match: '/editor/coffee/tastings' }
    ]
  },
  {
    label: 'editor.music',
    items: [{ label: 'editor.playlists', to: '/editor/music', match: '/editor/music' }]
  },
  {
    label: 'editor.travel',
    items: [{ label: 'editor.mapLocations', to: '/editor/travel', match: '/editor/travel' }]
  }
] as const

interface Entry {
  to: string
  match: string
}

function isActive(item: Entry): boolean {
  return route.path === item.match || route.path.startsWith(`${item.match}/`)
}
</script>
