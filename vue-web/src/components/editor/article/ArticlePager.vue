<template>
  <div class="flex flex-wrap items-center gap-3">
    <span class="text-sm font-medium text-foreground/60">{{ summary }}</span>

    <div class="ml-auto flex flex-wrap items-center gap-1">
      <EditorButton
        size="sm"
        variant="secondary"
        :disabled="page <= 1"
        @click="$emit('change', page - 1)"
      >
        {{ $t('common.previous') }}
      </EditorButton>

      <EditorButton
        v-for="value in pages"
        :key="value"
        size="sm"
        variant="secondary"
        :class="value === page ? 'bg-shade' : ''"
        :disabled="value === page"
        @click="$emit('change', value)"
      >
        {{ value }}
      </EditorButton>

      <EditorButton
        size="sm"
        variant="secondary"
        :disabled="page >= last_page"
        @click="$emit('change', page + 1)"
      >
        {{ $t('common.next') }}
      </EditorButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import EditorButton from '@/components/editor/ui/EditorButton.vue'

interface Props {
  page: number
  page_size: number
  total: number
  total_pages: number
}

const $props = defineProps<Props>()

defineEmits<{ change: [page: number] }>()

const { t } = useI18n()

const last_page = computed(() => Math.max($props.total_pages, 1))
const start = computed(() => ($props.total === 0 ? 0 : ($props.page - 1) * $props.page_size + 1))
const end = computed(() => Math.min($props.page * $props.page_size, $props.total))

const summary = computed(() =>
  t('article.editor.pageSummary', {
    start: start.value,
    end: end.value,
    total: $props.total
  })
)

const pages = computed(() => {
  const last = last_page.value
  const first = Math.max(1, Math.min($props.page - 2, last - 4))
  const final = Math.min(last, first + 4)
  const values: number[] = []

  for (let value = first; value <= final; value += 1) values.push(value)

  return values
})
</script>
