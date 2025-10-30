import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  base: './',
  
  // 关键：添加路径解析
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      '@/assets': resolve(__dirname, 'src/assets')
    }
  },
  
  build: {
    outDir: 'dist',
    target: 'es2015',
    
    rollupOptions: {
      output: {
        format: 'iife',
        chunkFileNames: '[name]-[hash].js',
        entryFileNames: '[name]-[hash].js',
        assetFileNames: 'img/[name]-[hash][extname]'
      }
    },
    
    modulePreload: false,
    minify: 'esbuild'
  }
})