import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    // dev では /api を api(:8080)へプロキシ。本番は web(Caddy)が同じ /api を
    // api へ振り分けるため、相対パスのまま単一オリジンで叩ける(ADR-0015)。
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
})
