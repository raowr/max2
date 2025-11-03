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
        // 修改这里：使用函数来区分不同类型的资源
        assetFileNames: (assetInfo) => {
          // 对于音频文件，输出到 music 目录
          if (assetInfo.name && /\.(mp3|wav|ogg)$/.test(assetInfo.name)) {
            return 'music/[name]-[hash][extname]'
          }
          // 其他资源仍然输出到 img 目录
          return 'img/[name]-[hash][extname]'
        }
      }
    },
    
    modulePreload: false,
    minify: 'esbuild'
  }
})