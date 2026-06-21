import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    // dev では api(:8080)へプロキシし、本番の単一オリジンと同じ相対パスで叩く
    proxy: {
      '/healthz': 'http://localhost:8080',
      '/api': 'http://localhost:8080',
    },
  },
})
