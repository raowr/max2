import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  // 重要：配置构建选项以适应移动端
  build: {
    outDir: 'dist',
    assetsDir: 'assets',
    // 确保资源路径正确
    assetsInlineLimit: 4096,
  },
  // 开发服务器配置
  server: {
    host: '0.0.0.0', // 允许外部访问（重要！）
    port: 5175
  }
})