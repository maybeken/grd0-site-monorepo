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
      <EditorButton
        variant="icon"
        :label="$t('editor.closeMenu')"
        class="lg:hidden"
        @click="$emit('close')"
      >
        <Icon icon="mynaui:x"></Icon>
      </EditorButton>
    </div>

    <nav class="flex flex-1 flex-col gap-6 overflow-y-auto px-3 pb-6">
      <div v-for="group in groups" :key="group.label">
        <p
          class="px-2 pb-2 text-sm font-medium uppercase tracking-widest text-foreground/60"
        >
          {{ $t(group.label) }}
        </p>
        <RouterLink
          v-for="item in group.items"
          :key="item.to"
          :to="item.to"
          class="block rounded px-2 py-2 text-lg font-medium transition-colors"
          :class="
            isActive(item)
              ? 'bg-accent font-semibold text-foreground'
              : 'text-foreground/60 hover:bg-accent/10 hover:text-foreground'
          "
          :aria-current="isActive(item) ? 'page' : undefined"
        >
          {{ $t(item.label) }}
        </RouterLink>
      </div>
    </nav>

    <div class="border-t border-accent/20 p-3">
      <RouterLink
        to="/"
        class="block rounded px-2 py-2 text-lg font-medium text-foreground/60 hover:bg-accent/10"
      >
        {{ $t('editor.viewSite') }}
      </RouterLink>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

import EditorButton from '@/components/editor/ui/EditorButton.vue'
import Icon from '@/components/shared/Icon.vue'

defineProps<{
  open: boolean
}>()

defineEmits<{
  close: []
}>()

const route = useRoute()

interface Entry {
  label: string
  to: string
  match: string
}

interface Group {
  label: string
  items: Entry[]
}

const groups: Group[] = [
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
    items: [
      { label: 'editor.songs', to: '/editor/music/songs', match: '/editor/music/songs' },
      { label: 'editor.playlists', to: '/editor/music', match: '/editor/music' }
    ]
  },
  {
    label: 'editor.travel',
    items: [{ label: 'editor.mapLocations', to: '/editor/travel', match: '/editor/travel' }]
  }
]

const active_match = computed(() => {
  const matches = groups
    .flatMap((group) => group.items)
    .filter((item) => route.path === item.match || route.path.startsWith(`${item.match}/`))

  return matches.sort((left, right) => right.match.length - left.match.length)[0]?.match ?? ''
})

function isActive(item: Entry): boolean {
  return active_match.value === item.match
}
</script>
