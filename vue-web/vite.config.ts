import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import { VitePWA } from 'vite-plugin-pwa'
import Components from 'unplugin-vue-components/vite'
import gettextYamlI18nPlugin from './vite-plugin-gettext-yaml-i18n'

import type { Plugin } from 'vite'

// ONNX Runtime references its ~27 MB wasm binary from the worker chunk, which
// Vite then emits. The runtime loads that binary from the asset CDN instead,
// so the copy is dropped: Cloudflare Pages rejects files above 25 MiB.
function dropOrtWasm(): Plugin {
  return {
    name: 'drop-ort-wasm',
    apply: 'build',
    generateBundle(_options, bundle) {
      for (const file of Object.keys(bundle)) {
        if (/ort-wasm-.*\.(wasm|mjs)$/.test(file)) delete bundle[file]
      }
    }
  }
}

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    gettextYamlI18nPlugin(),
    vue(),
    tailwindcss(),
    dropOrtWasm(),
    Components({
      excludeNames: ['Map'],
    }),
    VitePWA({
      registerType: 'autoUpdate',
      includeAssets: ['/assets/pwa-192x192.png', '/assets/pwa-512x512.png', '/assets/pwa-maskable-192x192.png', '/assets/pwa-maskable-512x512.png'],
      workbox: {
        globIgnores: [
          '**/cutout.worker-*.js',
          '**/cutout-runtime-*.js',
          '**/ort-wasm-*.wasm',
          '**/ort-wasm-*.mjs'
        ]
      },
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
  build: {
    rollupOptions: {
      output: {
        manualChunks: (id: string) => {
          if (
            id.includes('@huggingface/transformers') ||
            id.includes('onnxruntime-web') ||
            id.includes('onnxruntime-common')
          ) {
            return 'cutout-runtime'
          }
        }
      }
    }
  },
})
