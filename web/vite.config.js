import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import path from 'path'; // 新增这行导入

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src') // 映射 @ 为 src 目录
    },
  },
    build: {
      outDir: path.resolve(__dirname, '../max2/app/web/resource/dist'),
      rollupOptions: {
        input: path.resolve(__dirname, 'index.html') // 明确指定入口文件
      }
    },
  
})
