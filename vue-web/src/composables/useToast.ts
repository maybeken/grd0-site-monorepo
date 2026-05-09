import { reactive } from 'vue'

export interface Toast {
  id: number
  type: 'success' | 'error' | 'info'
  message: string
  duration: number
}

const state = reactive<{
  toasts: Toast[]
  nextId: number
}>({
  toasts: [],
  nextId: 0,
})

export function useToast() {
  const addToast = (type: Toast['type'], message: string, duration: number) => {
    const id = state.nextId++
    state.toasts.push({ id, type, message, duration })

    if (duration > 0) {
      setTimeout(() => {
        removeToast(id)
      }, duration)
    }
  }

  const removeToast = (id: number) => {
    const index = state.toasts.findIndex((t) => t.id === id)
    if (index !== -1) {
      state.toasts.splice(index, 1)
    }
  }

  const success = (message: string, duration = 4000) => {
    addToast('success', message, duration)
  }

  const error = (message: string, duration = 6000) => {
    addToast('error', message, duration)
  }

  const info = (message: string, duration = 4000) => {
    addToast('info', message, duration)
  }

  return {
    state,
    success,
    error,
    info,
    removeToast,
  }
}
