import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    host: true,
    watch: {
      usePolling: true,    
      interval: 100,
    },
    proxy: {
      '/api': {
        target: 'http://backend:8080',
        changeOrigin: true,
      },
      '/login': {
        target: 'http://backend:8080',
        changeOrigin: true,
      },
      '/refresh': {
        target: 'http://backend:8080',
        changeOrigin: true,
      },
    },
  },
})
