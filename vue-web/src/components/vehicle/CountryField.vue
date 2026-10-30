<template>
  <div class="relative">
    <div
      class="flex items-center gap-2 rounded border border-foreground/40 bg-background px-2.5 py-1.5 transition-colors focus-within:border-foreground/80 focus-within:outline focus-within:outline-2 focus-within:outline-offset-2 focus-within:outline-foreground/60"
    >
      <Icon
        v-if="preview_flag"
        :icon="preview_flag"
        dotify="false"
        height="1.25rem"
        class="shrink-0"
      ></Icon>

      <input
        :id="id"
        ref="input"
        :value="query"
        type="text"
        role="combobox"
        autocomplete="off"
        :aria-expanded="open"
        :aria-controls="listId"
        aria-autocomplete="list"
        :aria-activedescendant="active_option_id"
        :placeholder="placeholder"
        class="w-full bg-transparent text-lg font-medium outline-none placeholder:text-foreground/60"
        @input="onInput"
        @focus="onFocus"
        @keydown.down.prevent="move(1)"
        @keydown.up.prevent="move(-1)"
        @keydown.enter.prevent="commit"
        @keydown.esc.prevent="closeAndRevert"
        @blur="onBlur"
      />

      <button
        v-if="query"
        type="button"
        class="shrink-0 rounded p-0.5 text-foreground/60 hover:text-foreground"
        :aria-label="$t('grandTour.editor.form.clearCountry')"
        @mousedown.prevent
        @click="clear"
      >
        <Icon icon="mynaui:x" height="1rem"></Icon>
      </button>
    </div>

    <ul
      v-if="open && suggestions.length"
      :id="listId"
      role="listbox"
      class="absolute z-30 max-h-64 w-full overflow-y-auto rounded border border-accent/40 bg-background shadow-lg"
      :class="placement === 'top' ? 'bottom-full mb-1' : 'top-full mt-1'"
    >
      <li
        v-for="(option, index) in suggestions"
        :id="optionId(option)"
        :key="option.code"
        role="option"
        :aria-selected="index === active"
        class="flex cursor-pointer items-center gap-2 px-2.5 py-1.5 text-lg"
        :class="index === active ? 'bg-accent/30' : 'hover:bg-accent/20'"
        @mousedown.prevent="select(option)"
        @mouseenter="active = index"
      >
        <Icon
          :icon="`circle-flags:${option.code.toLowerCase()}`"
          dotify="false"
          height="1.25rem"
          class="shrink-0"
        ></Icon>
        <span class="truncate">{{ option.name }}</span>
        <span class="ml-auto text-sm text-foreground/60">{{ option.code }}</span>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'

import Icon from '@/components/shared/Icon.vue'
import { hasFlag, loadFlagCollection } from '@/helpers/icons'
import { countryOptions } from '@/helpers/vehicle'

import type { CountryOption } from '@/helpers/vehicle'

const props = withDefaults(
  defineProps<{
    code?: string | null
    country?: string | null
    id?: string
    placeholder?: string
  }>(),
  {
    code: null,
    country: null,
    id: 'vehicle-country',
    placeholder: ''
  }
)

const emit = defineEmits<{
  select: [selection: { code: string | null; country: string | null }]
}>()

const options = ref<CountryOption[]>([])
const input = ref<HTMLInputElement | null>(null)
const open = ref(false)
const active = ref(0)
const query = ref('')
const placement = ref<'bottom' | 'top'>('bottom')

const listId = `${props.id}-countries`

function optionId(option: CountryOption): string {
  return `${props.id}-country-${option.code}`
}

const selected_name = computed(() => {
  const code = (props.code ?? '').toUpperCase()
  const option = options.value.find((item) => item.code === code)

  return option?.name ?? props.country ?? ''
})

const suggestions = computed(() => {
  const term = query.value.trim().toLowerCase()

  if (!term) return options.value.slice(0, 50)

  return options.value
    .filter(
      (option) =>
        option.name.toLowerCase().includes(term) ||
        option.code.toLowerCase().startsWith(term)
    )
    .slice(0, 50)
})

const active_option = computed(() =>
  open.value ? (suggestions.value[active.value] ?? null) : null
)

const preview_flag = computed(() => {
  const code = (active_option.value?.code ?? props.code ?? '').toLowerCase()

  return code.length === 2 && hasFlag(code) ? `circle-flags:${code}` : null
})

const active_option_id = computed(() =>
  active_option.value ? optionId(active_option.value) : undefined
)

watch(
  [() => props.code, () => props.country],
  () => {
    if (!open.value) query.value = selected_name.value
  },
  { immediate: true }
)

function updatePlacement(): void {
  const rect = input.value?.getBoundingClientRect()

  if (!rect) return

  placement.value = rect.bottom + 280 > window.innerHeight ? 'top' : 'bottom'
}

function onFocus(): void {
  updatePlacement()
  open.value = true
  active.value = 0
}

function onInput(event: Event): void {
  query.value = (event.target as HTMLInputElement).value
  updatePlacement()
  open.value = true
  active.value = 0
}

function move(delta: number): void {
  open.value = true

  const count = suggestions.value.length

  if (count === 0) return

  active.value = (active.value + delta + count) % count
}

function select(option: CountryOption): void {
  query.value = option.name
  open.value = false
  emit('select', { code: option.code, country: option.name })
}

function commit(): void {
  const option = suggestions.value[active.value]

  if (option) {
    select(option)
    return
  }

  closeAndRevert()
}

function clear(): void {
  query.value = ''
  open.value = false
  emit('select', { code: null, country: null })
  input.value?.focus()
}

function closeAndRevert(): void {
  open.value = false
  query.value = selected_name.value
}

function onBlur(): void {
  if (options.value.length === 0) {
    const text = query.value.trim()

    emit('select', { code: null, country: text || null })
    return
  }

  closeAndRevert()
}

onMounted(async () => {
  try {
    await loadFlagCollection()
  } catch {
    options.value = []
  }

  options.value = countryOptions()

  if (options.value.length > 0) query.value = selected_name.value
})
</script>
