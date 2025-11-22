import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), tailwindcss()],
  server: {
    port: 5173,
    host: '0.0.0.0',
    watch: {
      usePolling: true,  // Enable polling for file changes in Docker on Windows
      interval: 100      // Check for changes every 100ms
    },
    proxy: {
      '/api': {
        target: 'http://nginx',  // ← Direct container-to-container communication
        changeOrigin: true,
        secure: false,
      }
    }
  }
})
