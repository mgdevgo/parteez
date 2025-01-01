import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from "path"

// https://vite.dev/config/
export default defineConfig({
  root: path.join(__dirname, "src"),
  build: {
    outDir: path.resolve(__dirname, "build"),
    emptyOutDir: true,
    rollupOptions: {
      input: {
        index: path.resolve(__dirname, 'src/index.html'),
      }
    }
  },
  plugins: [react()],
})
