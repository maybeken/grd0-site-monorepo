<template>
  <Teleport to="body">
    <div class="fixed bottom-4 right-4 z-50 flex flex-col gap-2 pointer-events-none">
      <TransitionGroup name="toast" tag="div">
        <div
          v-for="toast in state.toasts"
          :key="toast.id"
          class="pointer-events-auto px-4 py-3 rounded-xl shadow-lg max-w-sm flex items-start gap-3 border backdrop-blur-sm"
          :class="{
            'bg-green-900/80 border-green-700 text-green-100': toast.type === 'success',
            'bg-red-900/80 border-red-700 text-red-100': toast.type === 'error',
            'bg-blue-900/80 border-blue-700 text-blue-100': toast.type === 'info',
          }"
        >
          <span class="text-sm flex-1">{{ toast.message }}</span>
          <button
            class="text-current opacity-50 hover:opacity-100 transition-opacity flex-shrink-0"
            @click="removeToast(toast.id)"
          >
            ✕
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { useToast } from '@/composables/useToast'

const { state, removeToast } = useToast()
</script>

<style scoped>
.toast-enter-active {
  transition: all 0.3s ease-out;
}
.toast-leave-active {
  transition: all 0.2s ease-in;
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}
.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>
