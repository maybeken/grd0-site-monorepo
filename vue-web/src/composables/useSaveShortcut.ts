import { onMounted, onUnmounted } from 'vue'

export interface SaveShortcutOptions {
  enabled?: () => boolean
}

export function useSaveShortcut(
  save: () => void | Promise<void>,
  options: SaveShortcutOptions = {}
): void {
  function on_keydown(event: KeyboardEvent): void {
    if (!(event.metaKey || event.ctrlKey) || event.altKey || event.shiftKey) return
    if (event.key.toLowerCase() !== 's') return
    if (options.enabled && !options.enabled()) return

    event.preventDefault()
    void save()
  }

  onMounted(() => window.addEventListener('keydown', on_keydown))
  onUnmounted(() => window.removeEventListener('keydown', on_keydown))
}
