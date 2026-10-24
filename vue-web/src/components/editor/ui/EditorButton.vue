<template>
  <component
    :is="to ? RouterLink : 'button'"
    :to="to || undefined"
    :type="to ? undefined : type"
    :disabled="to ? undefined : disabled || loading"
    :form="to ? undefined : form"
    :aria-label="label || undefined"
    :title="label || undefined"
    :aria-busy="loading ? 'true' : undefined"
    class="inline-flex items-center justify-center gap-2 rounded font-medium transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-foreground/60 disabled:cursor-not-allowed disabled:opacity-50"
    :class="[variant_class, size_class]"
  >
    <Icon v-if="loading" icon="mynaui:loader" height="1rem" class="animate-spin"></Icon>
    <slot></slot>
  </component>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'

import Icon from '@/components/shared/Icon.vue'

import type { RouteLocationRaw } from 'vue-router'

type Variant = 'primary' | 'secondary' | 'destructive' | 'icon'

const props = withDefaults(
  defineProps<{
    variant?: Variant
    size?: 'md' | 'sm'
    type?: 'button' | 'submit' | 'reset'
    disabled?: boolean
    loading?: boolean
    label?: string
    to?: RouteLocationRaw
    form?: string
  }>(),
  {
    variant: 'secondary',
    size: 'md',
    type: 'button',
    disabled: false,
    loading: false,
    label: '',
    to: undefined,
    form: undefined
  }
)

if (import.meta.env.DEV && props.variant === 'icon' && !props.label) {
  console.warn('[EditorButton] The icon variant requires a visible label for accessibility.')
}

const variant_class = computed(() => {
  if (props.variant === 'primary') {
    return 'bg-accent px-3 py-1.5 text-lg font-semibold text-foreground hover:brightness-110'
  }

  if (props.variant === 'destructive') {
    return 'border border-red-500/40 px-3 py-1.5 text-lg text-red-400 hover:bg-red-500/10'
  }

  if (props.variant === 'icon') {
    return 'rounded-full border border-dotted p-1 hover:bg-accent/10'
  }

  return 'border border-accent/40 px-3 py-1.5 text-lg hover:bg-accent/10'
})

const size_class = computed(() =>
  props.size === 'sm' && props.variant !== 'icon' ? 'px-2 py-1' : ''
)
</script>
