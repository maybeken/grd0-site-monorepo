import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

import Components from 'unplugin-vue-components/vite'

function manualChunks (id: any) {
  if (id.includes('/vue/') || id.includes('/vue-router/') || id.includes('/@vue/')) {
    return 'vue';
  } else if (id.includes('/alova/')) {
    return 'alova';
  } else if (id.includes('/rehype/') || id.includes('/remark/') || id.includes('micromark') || id.includes('markdown') || id.includes('/unified/')) {
    return 'md';
  } else if (id.includes('/icons/')) {
    return 'icon';
  } else if (id.includes('/ol/') || id.includes('/vue3-openlayers/')) {
    return 'openlayers';
  } else if (id.includes('node_modules')) {
    return 'vendor';
  } else {
    return 'main';
  }
}

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    Components(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks
      }
    }
  }
})
