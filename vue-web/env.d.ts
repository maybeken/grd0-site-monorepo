/// <reference types='vite/client' />

declare module '*.yaml' {
  const data: Record<string, any>
  export default data
}

declare module '@intlify/unplugin-vue-i18n/vite' {
  import type { Plugin } from 'vite'
  interface PluginOptions {
    include?: string | string[]
  }
  export default function vueI18n(options?: PluginOptions): Plugin
}
