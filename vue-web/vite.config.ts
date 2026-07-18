import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import { VitePWA } from 'vite-plugin-pwa'
import Components from 'unplugin-vue-components/vite'
import gettextYamlI18nPlugin from './vite-plugin-gettext-yaml-i18n'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    gettextYamlI18nPlugin(),
    vue(),
    tailwindcss(),
    Components({
      excludeNames: ['Map'],
    }),
    VitePWA({
      registerType: 'autoUpdate',
      includeAssets: ['/assets/pwa-192x192.png', '/assets/pwa-512x512.png', '/assets/pwa-maskable-192x192.png', '/assets/pwa-maskable-512x512.png'],
      manifest: {
        name: 'IDEA GROUND ZERO',
        short_name: 'grd0.net',
        description: 'Blog post, photos, and useful tools from grd0.net',
        theme_color: '#0A0A0A',
        icons: [
          {
            src: '/assets/pwa-192x192.png',
            sizes: '192x192',
            type: 'image/png'
          },
          {
            src: '/assets/pwa-512x512.png',
            sizes: '512x512',
            type: 'image/png'
          },
          {
            src: '/assets/pwa-512x512.png',
            sizes: '512x512',
            type: 'image/png',
            purpose: 'any'
          },
          {
            src: '/assets/pwa-512x512.png',
            sizes: '512x512',
            type: 'image/png',
            purpose: 'maskable'
          }
        ],
      },
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
})
