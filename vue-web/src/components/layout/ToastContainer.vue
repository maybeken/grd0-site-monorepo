<template>
  <div
    class="fixed bottom-4 left-4 right-4 sm:left-auto sm:right-4 z-[100] flex flex-col gap-2 pointer-events-none sm:w-96"
    aria-live="polite"
    aria-atomic="false"
  >
    <TransitionGroup name="toast">
      <div
        v-for="toast in toasts"
        :key="toast.id"
        role="status"
        class="pointer-events-auto flex items-start gap-3 rounded-xl border bg-background/95 px-4 py-2 text-sm shadow-lg backdrop-blur"
        :class="type_classes[toast.type]"
      >
        <span class="grow break-words">{{ toast.message }}</span>
        <button
          type="button"
          class="shrink-0 text-secondary hover:text-foreground"
          :aria-label="$t('common.dismiss')"
          :title="$t('common.dismiss')"
          @click="dismiss(toast.id)"
        >
          <Icon icon="mynaui:x" height="1rem"></Icon>
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup lang="ts">
import { useToast } from '@/composables/useToast'

import type { ToastType } from '@/composables/useToast'

const { toasts, dismiss } = useToast()

const type_classes: Record<ToastType, string> = {
  success: 'border-emerald-500',
  error: 'border-red-400',
  info: 'border-accent'
}
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition:
    opacity 0.2s ease,
    transform 0.2s ease;
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(0.5rem);
}
</style>
