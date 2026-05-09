<template>
  <Teleport to="body">
    <div
      v-if="state.isOpen"
      class="fixed inset-0 z-50 flex items-center justify-center"
      role="dialog"
      aria-modal="true"
      :aria-label="state.title"
    >
      <div class="fixed inset-0 bg-black/50 backdrop-blur-sm" @click="close"></div>

      <div
        ref="modalRef"
        class="relative bg-gray-900 border border-gray-700 rounded-2xl p-6 w-full max-w-md mx-4 shadow-2xl"
        tabindex="-1"
      >
        <h2 class="text-lg font-semibold text-white mb-2">{{ state.title }}</h2>
        <p class="text-gray-300 mb-4">{{ state.message }}</p>

        <input
          v-if="state.type === 'prompt'"
          ref="inputRef"
          v-model="state.inputValue"
          :placeholder="state.inputPlaceholder"
          class="w-full bg-gray-800 border border-gray-600 rounded-xl px-4 py-2 text-white placeholder-gray-500 mb-4 focus:outline-none focus:border-gray-400"
          @keydown.enter="submit"
        />

        <div class="flex justify-end gap-3">
          <button
            class="px-4 py-2 rounded-xl border border-gray-600 text-gray-300 hover:bg-gray-800 transition-colors"
            @click="close"
          >
            Cancel
          </button>
          <button
            class="px-4 py-2 rounded-xl bg-white text-black hover:bg-gray-200 transition-colors"
            @click="submit"
          >
            {{ state.type === 'prompt' ? 'OK' : 'Confirm' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { useModal } from '@/composables/useModal'

const { state, close, submit } = useModal()

const modalRef = ref<HTMLElement | null>(null)
const inputRef = ref<HTMLInputElement | null>(null)

const focusableSelector =
  'a[href], button:not([disabled]), textarea:not([disabled]), input:not([disabled]), select:not([disabled]), [tabindex]:not([tabindex="-1"])'

function trapFocus(e: KeyboardEvent) {
  if (e.key !== 'Tab' || !modalRef.value) return

  const focusable = modalRef.value.querySelectorAll<HTMLElement>(focusableSelector)
  if (focusable.length === 0) return

  const first = focusable[0]!
  const last = focusable[focusable.length - 1]!

  if (e.shiftKey && document.activeElement === first) {
    e.preventDefault()
    last.focus()
  } else if (!e.shiftKey && document.activeElement === last) {
    e.preventDefault()
    first.focus()
  }
}

function handleEscape(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    close()
  }
}

watch(
  () => state.isOpen,
  async (isOpen) => {
    if (isOpen) {
      document.addEventListener('keydown', handleEscape)
      document.addEventListener('keydown', trapFocus)
      await nextTick()
      modalRef.value?.focus()

      if (state.type === 'prompt') {
        await nextTick()
        inputRef.value?.focus()
      }
    } else {
      document.removeEventListener('keydown', handleEscape)
      document.removeEventListener('keydown', trapFocus)
    }
  }
)
</script>
