import { createI18n } from 'vue-i18n'
import en from './locales/en.yaml'

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en
  }
})

export const t = (key: string, values?: Record<string, unknown>): string =>
  (i18n.global as unknown as { t: (key: string, values?: Record<string, unknown>) => string }).t(
    key,
    values
  )

export default i18n
