import { useI18n } from 'vue-i18n'

export function useConfirmDiscard(): (dirty: boolean, messageKey?: string) => boolean {
  const { t } = useI18n()

  return (dirty: boolean, messageKey = 'common.confirmDiscard'): boolean => {
    if (!dirty) return true

    return window.confirm(t(messageKey))
  }
}
