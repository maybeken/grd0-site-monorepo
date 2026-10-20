import { reactive } from 'vue'

export type ToastType = 'success' | 'error' | 'info'

export interface ToastMessage {
  id: number
  message: string
  type: ToastType
}

const toasts = reactive<ToastMessage[]>([])

let next_id = 0

export function dismissToast(id: number): void {
  const index = toasts.findIndex((toast) => toast.id === id)
  if (index >= 0) toasts.splice(index, 1)
}

function showToast(message: string, type: ToastType, duration: number): number {
  const id = ++next_id
  toasts.push({ id, message, type })

  if (toasts.length > 5) toasts.shift()

  if (duration > 0) {
    window.setTimeout(() => dismissToast(id), duration)
  }

  return id
}

export function useToast() {
  return {
    toasts,
    dismiss: dismissToast,
    show: (message: string, type: ToastType = 'info', duration = 4000) =>
      showToast(message, type, duration),
    success: (message: string, duration = 4000) => showToast(message, 'success', duration),
    error: (message: string, duration = 6000) => showToast(message, 'error', duration),
    info: (message: string, duration = 4000) => showToast(message, 'info', duration)
  }
}
