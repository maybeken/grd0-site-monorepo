import { onUnmounted } from 'vue'
import { onBeforeRouteLeave, onBeforeRouteUpdate } from 'vue-router'
import { useI18n } from 'vue-i18n'

import type { Ref } from 'vue'

export function useUnsavedGuard(
  dirty: Ref<boolean>,
  messageKey = 'common.confirmDiscard'
): void {
  const { t } = useI18n()

  function confirmDiscard(): boolean {
    if (!dirty.value) return true

    return window.confirm(t(messageKey))
  }

  function onBeforeUnload(event: BeforeUnloadEvent): void {
    if (!dirty.value) return

    event.preventDefault()
    event.returnValue = ''
  }

  onBeforeRouteLeave(confirmDiscard)
  onBeforeRouteUpdate(confirmDiscard)

  window.addEventListener('beforeunload', onBeforeUnload)
  onUnmounted(() => window.removeEventListener('beforeunload', onBeforeUnload))
}
