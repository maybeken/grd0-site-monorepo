<template>
  <div class="flex flex-wrap items-center gap-3 text-sm">
    <span class="text-secondary">{{ summary }}</span>

    <div class="ml-auto flex flex-wrap items-center gap-1">
      <button
        class="rounded-full border border-dotted px-3 py-1 disabled:opacity-40"
        :disabled="page <= 1"
        @click="$emit('change', page - 1)"
      >
        {{ $t('common.previous') }}
      </button>

      <button
        v-for="value in pages"
        :key="value"
        class="rounded-full border border-dotted px-3 py-1 disabled:opacity-40"
        :class="value === page ? 'bg-shade' : ''"
        :disabled="value === page"
        @click="$emit('change', value)"
      >
        {{ value }}
      </button>

      <button
        class="rounded-full border border-dotted px-3 py-1 disabled:opacity-40"
        :disabled="page >= last_page"
        @click="$emit('change', page + 1)"
      >
        {{ $t('common.next') }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

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
