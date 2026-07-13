import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    outDir: '../agent/static',
    assetsDir: 'assets',
    sourcemap: true,
  },
})
