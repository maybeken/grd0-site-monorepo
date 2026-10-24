<template>
  <div ref="container" class="relative">
    <button
      type="button"
      class="flex w-full items-center gap-2 rounded border border-foreground/40 bg-background px-2.5 py-1.5 text-left text-lg font-medium hover:bg-accent/5"
      :aria-expanded="open"
      aria-haspopup="listbox"
      @click="toggle"
    >
      <span :class="modelValue ? 'font-mono text-lg' : 'text-foreground/60'">
        {{ modelValue || $t('map.editor.iconSelect') }}
      </span>
      <Icon icon="mynaui:chevron-down" class="ml-auto shrink-0" height="1rem" />
    </button>

    <div
      v-if="open"
      class="absolute inset-x-0 z-30 mt-1 flex max-h-96 flex-col gap-2 rounded-xl border border-accent/40 bg-background p-2 shadow-lg"
      @keydown.esc.stop="open = false"
    >
      <div class="relative">
        <Icon
          icon="mynaui:search"
          class="pointer-events-none absolute top-1/2 left-2 -translate-y-1/2 text-foreground/60"
          height="1rem"
        />
        <input
          ref="search_input"
          v-model="query"
          type="search"
          class="w-full rounded border border-foreground/40 bg-background py-1.5 pr-2 pl-7 text-lg font-medium placeholder:text-foreground/60 focus:border-foreground/80 focus:outline focus:outline-2 focus:outline-offset-2 focus:outline-foreground/60"
          :placeholder="$t('map.editor.iconSearchPlaceholder')"
          @keydown.enter.prevent="selectFirst"
        />
      </div>

      <p v-if="loading" class="py-4 text-center text-sm font-medium text-foreground/60">
        {{ $t('common.loading') }}
      </p>
      <p v-else-if="failed" class="py-4 text-center text-sm font-medium text-red-400">
        {{ $t('map.editor.iconLoadFailed') }}
      </p>
      <p v-else-if="!results.length" class="py-4 text-center text-sm font-medium text-foreground/60">
        {{ $t('map.editor.iconNoResults') }}
      </p>
      <div
        v-else
        class="grid min-h-0 flex-1 grid-cols-6 gap-1 overflow-y-auto sm:grid-cols-8"
        @scroll="onGridScroll"
      >
        <button
          v-for="name in results"
          :key="name"
          type="button"
          class="flex aspect-square items-center justify-center rounded border text-foreground/80 hover:bg-accent/10 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-foreground/60"
          :class="
            name === modelValue
              ? 'border-accent bg-accent/15 text-accent'
              : 'border-transparent'
          "
          :title="displayName(name)"
          :aria-label="displayName(name)"
          @click="select(name)"
        >
          <Icon :icon="name" height="1.5rem" />
        </button>
      </div>

      <div
        v-if="!loading && !failed"
        class="flex items-center justify-between gap-2 border-t border-accent/10 pt-2"
      >
        <span class="text-sm font-medium text-foreground/60">
          {{ $t('map.editor.iconShowing', { visible: results.length, total: resultsAll.length }) }}
        </span>
        <div class="flex shrink-0 items-center gap-2">
          <EditorButton
            v-if="hasMore"
            size="sm"
            variant="secondary"
            @click="loadMore"
          >
            {{ $t('common.loadMore') }}
          </EditorButton>
          <EditorButton
            v-if="modelValue"
            size="sm"
            variant="secondary"
            @click="clear"
          >
            {{ $t('map.editor.iconClear') }}
          </EditorButton>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'

import EditorButton from '@/components/editor/ui/EditorButton.vue'
import Icon from '@/components/shared/Icon.vue'
import { loadIconCollection } from '@/helpers/icons'

const PREFIX = 'mynaui'
const PAGE_SIZE = 120

const POPULAR_ICONS = [
  'map-pin',
  'map',
  'location',
  'compass',
  'globe',
  'mountain',
  'mountain-snow',
  'heart',
  'star',
  'home',
  'building',
  'camera',
  'plane',
  'car',
  'train',
  'ship',
  'flag',
  'coffee',
  'food',
  'tent',
  'beach',
  'sun',
  'moon',
  'tree',
  'water',
  'route',
  'road',
  'navigation',
  'anchor',
  'rocket'
]

interface Props {
  modelValue: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const container = ref<HTMLElement | null>(null)
const search_input = ref<HTMLInputElement | null>(null)
const open = ref(false)
const query = ref('')
const icons = ref<string[]>([])
const loading = ref(false)
const failed = ref(false)
const visible_count = ref(PAGE_SIZE)

let loaded = false

async function loadIcons(): Promise<void> {
  if (loaded || loading.value) return

  loading.value = true
  failed.value = false

  try {
    const data = await loadIconCollection()
    const prefix = typeof data.prefix === 'string' && data.prefix ? data.prefix : PREFIX
    const names = [...Object.keys(data.icons ?? {}), ...Object.keys(data.aliases ?? {})]

    icons.value = Array.from(
      new Set(
        names.map((name) => (name.includes(':') ? name : `${prefix}:${name}`))
      )
    )
    loaded = true
  } catch {
    failed.value = true
  } finally {
    loading.value = false
  }
}

async function toggle(): Promise<void> {
  open.value = !open.value

  if (!open.value) return

  if (!loaded) {
    void loadIcons()
  }

  await nextTick()
  search_input.value?.focus()
}

const popular_icons = computed(() => {
  const available = new Set(icons.value)
  const result: string[] = []

  for (const name of POPULAR_ICONS) {
    for (const candidate of [`${PREFIX}:${name}`, `${PREFIX}:${name}-solid`]) {
      if (available.has(candidate)) result.push(candidate)
    }
  }

  return result
})

const resultsAll = computed(() => {
  const needle = query.value.trim().toLowerCase()

  if (needle) {
    return icons.value.filter((name) => name.toLowerCase().includes(needle))
  }

  const ordered: string[] = []
  const seen = new Set<string>()

  function push(name: string): void {
    if (seen.has(name)) return

    seen.add(name)
    ordered.push(name)
  }

  if (props.modelValue) push(props.modelValue)
  for (const name of popular_icons.value) push(name)
  for (const name of icons.value) push(name)

  return ordered
})

const results = computed(() => resultsAll.value.slice(0, visible_count.value))
const hasMore = computed(() => resultsAll.value.length > visible_count.value)

watch(query, () => {
  visible_count.value = PAGE_SIZE
})

function loadMore(): void {
  if (!hasMore.value) return

  visible_count.value += PAGE_SIZE
}

function onGridScroll(event: Event): void {
  const element = event.target as HTMLElement

  if (element.scrollTop + element.clientHeight >= element.scrollHeight - 48) {
    loadMore()
  }
}

function displayName(name: string): string {
  return name.startsWith(`${PREFIX}:`) ? name.slice(PREFIX.length + 1) : name
}

function select(name: string): void {
  emit('update:modelValue', name)
  open.value = false
  query.value = ''
}

function selectFirst(): void {
  const first = results.value[0]

  if (first) select(first)
}

function clear(): void {
  emit('update:modelValue', '')
  open.value = false
  query.value = ''
}

function onPointerDown(event: PointerEvent): void {
  if (!open.value) return
  if (container.value && !container.value.contains(event.target as Node)) {
    open.value = false
  }
}

onMounted(() => document.addEventListener('pointerdown', onPointerDown))
onUnmounted(() => document.removeEventListener('pointerdown', onPointerDown))
</script>
