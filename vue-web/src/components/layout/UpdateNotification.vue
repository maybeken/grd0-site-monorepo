<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'

const POLL_INTERVAL = 5 * 60 * 1000

const updateAvailable = ref(false)
const dismissed = ref(false)
const clearing = ref(false)
let pollTimer: number | undefined
let checking = false

function getEntryPath(root: ParentNode): string | null {
  const script = root.querySelector<HTMLScriptElement>('script[type="module"][src]')
  if (!script?.src) return null

  return new URL(script.src, window.location.origin).pathname
}

async function checkForUpdate() {
  if (import.meta.env.DEV || checking || document.visibilityState === 'hidden') return

  checking = true
  try {
    const current = getEntryPath(document)
    if (!current) return

    const response = await fetch(`${import.meta.env.BASE_URL}index.html?t=${Date.now()}`, {
      cache: 'no-store'
    })
    if (!response.ok) return

    const deployed = getEntryPath(
      new DOMParser().parseFromString(await response.text(), 'text/html')
    )
    if (!deployed) return

    updateAvailable.value = current !== deployed
  } catch {
    // Offline or the request failed; retry on the next focus or poll.
  } finally {
    checking = false
  }
}

async function clearCacheAndReload() {
  clearing.value = true
  try {
    if ('serviceWorker' in navigator) {
      const registrations = await navigator.serviceWorker.getRegistrations()
      await Promise.all(registrations.map((registration) => registration.unregister()))
    }

    if ('caches' in window) {
      const keys = await caches.keys()
      await Promise.all(keys.map((key) => caches.delete(key)))
    }
  } finally {
    window.location.reload()
  }
}

function handleVisibilityChange() {
  if (document.visibilityState === 'visible') checkForUpdate()
}

onMounted(() => {
  if (import.meta.env.DEV) return

  checkForUpdate()
  window.addEventListener('focus', checkForUpdate)
  window.addEventListener('online', checkForUpdate)
  document.addEventListener('visibilitychange', handleVisibilityChange)
  pollTimer = window.setInterval(checkForUpdate, POLL_INTERVAL)
})

onBeforeUnmount(() => {
  window.removeEventListener('focus', checkForUpdate)
  window.removeEventListener('online', checkForUpdate)
  document.removeEventListener('visibilitychange', handleVisibilityChange)
  window.clearInterval(pollTimer)
})
</script>

<template>
  <div
    v-if="updateAvailable && !dismissed"
    class="fixed right-4 bottom-4 z-[100] max-w-sm rounded-xl border border-accent bg-shade p-4 shadow-lg"
    role="status"
    aria-live="polite"
  >
    <p class="text-sm">{{ $t('updateNotification.message') }}</p>
    <div class="mt-3 flex items-center justify-end gap-2">
      <button
        class="rounded-lg px-3 py-1.5 text-xs text-secondary hover:text-foreground"
        @click="dismissed = true"
      >
        {{ $t('updateNotification.dismiss') }}
      </button>
      <button
        class="rounded-lg bg-foreground px-3 py-1.5 text-xs font-bold text-background disabled:opacity-50"
        :disabled="clearing"
        @click="clearCacheAndReload"
      >
        {{ clearing ? $t('updateNotification.clearing') : $t('updateNotification.action') }}
      </button>
    </div>
  </div>
</template>
